package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/config"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/services"
	"github.com/go-fuego/fuego"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type GoogleAuthController struct {
	config      *config.Config
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

func NewGoogleAuthController(
	config *config.Config,
	oauthConfig *oauth2.Config,
	jwtService *services.JWTService,
	userService *services.UserService,
) *GoogleAuthController {

	return &GoogleAuthController{
		config:      config,
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

func (gac *GoogleAuthController) Login(ctx fuego.ContextNoBody) (any, error) {
	state := generateRandomState(16)
	authURL := gac.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// store the state in a secure cookie
	http.SetCookie(ctx.Response(), &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // https
	})

	// redirect the user to Google's OAuth 2.0 consent page
	// this will then redirect back to our callback endpoint
	return ctx.Redirect(307, authURL)
}

func (gac *GoogleAuthController) Logout(ctx fuego.ContextNoBody) (any, error) {
	// Clear the JWT cookie
	http.SetCookie(ctx.Response(), &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // https
		MaxAge:   -1,    // delete the cookie immediately
	})

	return ctx.Redirect(307, gac.config.FrontendURL)
}

func (gac *GoogleAuthController) WhoAmI(ctx fuego.ContextNoBody) (jwt.Claims, error) {

	cookie, err := ctx.Cookie("jwt")
	if err != nil {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized, no JWT cookie found"}
	}

	claims, err := gac.jwtService.ValidateJWT(cookie.Value)
	if err != nil {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized: Invalid JWT token"}
	}

	// this returns the claims as JSON, see claims in jwt_service
	return claims.Claims, nil
}

func (gac *GoogleAuthController) Callback(ctx fuego.ContextNoBody) (any, error) {

	cookie, err := ctx.Cookie("oauthstate")
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "State Cookie not found"}
	}

	// Compare our client defined state with the state returned by Google
	state := ctx.Request().FormValue("state")
	if state != cookie.Value {
		return nil, fuego.BadRequestError{Detail: "Invalid OAuth state"}
	}

	code := ctx.Request().FormValue("code")

	token, err := gac.oauthConfig.Exchange(ctx.Context(), code)
	if err != nil {
		return nil, fuego.InternalServerError{Detail: "Failed to exchange token: " + err.Error()}
	}

	userInfo, err := fetchGoogleUserInfo(token)
	if err != nil {
		return nil, fuego.InternalServerError{Detail: "Failed to fetch Google user info: " + err.Error()}
	}

	// Check if user exists, if not create a new user
	dbUser, err := gac.userService.GetUserByProvider(context.Background(), "google", userInfo.ID)
	firstLogin := false
	if err != nil {
		// User does not exist, create new user
		userModel, err := domain.NewUser(userInfo.Name, userInfo.Email)
		if err != nil {
			return nil, fuego.InternalServerError{Detail: "Failed to create user: " + err.Error()}
		}

		// Set provider info in model
		// userModel.Provider = "google"
		// userModel.ProviderID = userInfo.ID
		userModel.UpdateProviderInformation("google", userInfo.ID)
		userModel.AvatarURL = userInfo.Picture

		// Register the user
		_, err = gac.userService.RegisterUser(context.Background(), userModel)
		if err != nil {
			return nil, fuego.InternalServerError{Detail: "Failed to register user: " + err.Error()}
		}

		dbUser = userModel
		firstLogin = true
	} else {
		// Update last login for existing user
		err = gac.userService.UpdateLastLogin(context.Background(), dbUser.ID)
		if err != nil {
			// erm... how did this happen !
		}
	}

	jwtToken, err := gac.jwtService.GenerateJWT(*dbUser)
	if err != nil {
		return nil, fuego.InternalServerError{Detail: "Failed to generate JWT: " + err.Error()}
	}

	// Set the JWT token in a secure cookie
	http.SetCookie(ctx.Response(), &http.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // https
	})

	if firstLogin {
		// Set first login cookie
		http.SetCookie(ctx.Response(), &http.Cookie{
			Name:     "first_login",
			Value:    strconv.FormatBool(firstLogin),
			Path:     "/",
			HttpOnly: false,
			MaxAge: 60, //Expire in 60 seconds
		})
	}

	return ctx.Redirect(307, gac.config.FrontendURL)
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
