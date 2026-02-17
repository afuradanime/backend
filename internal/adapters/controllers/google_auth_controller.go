package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/services"
	"golang.org/x/oauth2"
)

type GoogleAuthController struct {
	oauthConfig *oauth2.Config
	jwtService  *services.JWTService
	userService *services.UserService
}

type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewGoogleAuthController(oauthConfig *oauth2.Config, jwtService *services.JWTService,
	userService *services.UserService) *GoogleAuthController {

	return &GoogleAuthController{
		oauthConfig: oauthConfig,
		jwtService:  jwtService,
		userService: userService,
	}
}

// random state generator for OAuth2 flow.
func generateRandomState(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (gac *GoogleAuthController) Login(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState(16)
	authURL := gac.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// store the state in a secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		HttpOnly: true,
		Secure:   false, // https
	})

	// redirect the user to Google's OAuth 2.0 consent page
	// this will then redirect back to our callback endpoint
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (gac *GoogleAuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // https
		MaxAge:   -1,    // delete the cookie immediately
	})

	http.Redirect(w, r, os.Getenv("FRONTEND_URL"), http.StatusTemporaryRedirect)
}

func (gac *GoogleAuthController) WhoAmI(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "Unauthorized, no JWT cookie found", http.StatusUnauthorized)
		return
	}

	claims, err := gac.jwtService.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid JWT token", http.StatusUnauthorized)
		return
	}

	// this returns the claims as JSON, see claims in jwt_service
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims.Claims)
}

func (gac *GoogleAuthController) Callback(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("oauthstate")
	if err != nil {
		http.Error(w, "State Cookie not found", http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")
	if state != cookie.Value {
		// Compare our client defined state with the state returned by Google
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")

	token, err := gac.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := fetchGoogleUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to fetch Google user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user exists, if not create a new user
	db_user, err := gac.userService.GetUserByProvider(context.Background(), "google", userInfo.ID)
	if err != nil {
		// User does not exist, create new user
		user_model, err := domain.NewUser(userInfo.Name, userInfo.Email)
		if err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Set provider info in model
		user_model.Provider = "google"
		user_model.ProviderID = userInfo.ID
		user_model.AvatarURL = userInfo.Picture

		// Register the user
		err = gac.userService.RegisterUser(context.Background(), user_model)
		if err != nil {
			http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		db_user = user_model
	} else {
		// Update last login for existing user
		err = gac.userService.UpdateLastLogin(context.Background(), db_user.ID)
		if err != nil {
			// erm... how did this happen !
		}
	}

	jwtToken, err := gac.jwtService.GenerateJWT(*db_user)
	if err != nil {
		http.Error(w, "Failed to generate JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the JWT token in a secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // https
	})

	http.Redirect(w, r, os.Getenv("FRONTEND_URL"), http.StatusTemporaryRedirect)
}

func fetchGoogleUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	// call mr google to get user info via their user endpoint
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	decodeErr := json.NewDecoder(res.Body).Decode(&userInfo)
	if decodeErr != nil {
		return nil, errors.New("Failed decoding the Google user info response: " + decodeErr.Error())
	}

	res.Body.Close()

	return &userInfo, nil
}
