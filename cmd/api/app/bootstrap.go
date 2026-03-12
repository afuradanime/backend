package app

import (
	"context"
	"strconv"
	"time"

	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Bootstrap(m *mongo.Database) {

	// Nuke EVERYTHING
	m.Drop(context.Background())

	// Create repositories
	userRepo := repositories.NewUserRepository(m)
	friendshipRepo := repositories.NewFriendshipRepository(m)

	// Bootstrap users and get their auto-generated IDs
	krayID, taikoID, testID := BootstrapUsers(context.Background(), userRepo)

	// Bootstrap friendships using the actual user IDs
	BootstrapFriendships(context.Background(), friendshipRepo, krayID, taikoID, testID)

	// Bootstrap translation suggestions
	descRepo := repositories.NewDescriptionTranslationRepository(m)
	BootstrapTranslations(context.Background(), descRepo)

	// Bootstrap user reports
	reportRepo := repositories.NewUserReportRepository(m)
	BootstrapReports(context.Background(), reportRepo)

	// Bootstrap a post conversation
	postRepo := repositories.NewPostRepository(m)
	BootstrapPostConversation(context.Background(), postRepo, krayID, taikoID, testID)

	// Bootstrap animelist
	animeListRepo := repositories.NewAnimeListRepository(m)
	BootstrapAnimeList(context.Background(), animeListRepo, krayID)

	// Database Indices
	BootstrapIndices(context.Background(), m)
}

func BootstrapUsers(ctx context.Context, userRepo *repositories.UserRepository) (krayID, taikoID, testID int) {

	userKray, err := domain.NewUser("KrayRui", "kray@afurada.anime")
	if err != nil {
		panic(err)
	}

	userKray.UpdateLocation("Porto")
	userKray.UpdateSocials([]string{
		"https://x.com/RuiIshigami",
		"https://github.com/Rui-San",
	})

	userKray.AddRole(value.UserRoleModerator)
	userKray.AddRole(value.UserRoleAdmin)

	userKray.UpdateAvatarURL("/pfps/d7dea5d3e09941f563dabf364b4db31cac63a5f1.png")
	userKray.CreatedAt = time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC)

	// Create user and get auto-generated ID
	_, err = userRepo.CreateUser(ctx, userKray)
	if err != nil {
		panic(err)
	}
	krayID = userKray.ID

	userTaiko, err := domain.NewUser("Sagiri719", "tiagobarrossao@gmail.com")
	if err != nil {
		panic(err)
	}

	userTaiko.UpdateProviderInformation("google", "111642040238696442904")
	userTaiko.UpdateLocation("Porto")
	userTaiko.UpdateSocials([]string{
		"https://x.com/Sagiri719",
		"https://github.com/Sagiri721",
	})

	userTaiko.AddRole(value.UserRoleModerator)
	userTaiko.AddRole(value.UserRoleAdmin)

	userTaiko.UpdateAvatarURL("/pfps/e59084c01caf44df3c240a3c78009d080ea02556.png")
	userTaiko.RewardBadge(value.UserBadgeSuperMegaIllyaFan)

	userTaiko.CreatedAt = time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC)

	// Create user and get auto-generated ID
	_, err = userRepo.CreateUser(ctx, userTaiko)
	if err != nil {
		panic(err)
	}
	taikoID = userTaiko.ID

	userTest, err := domain.NewUser("Afuradanime", "afuradanime@gmail.com")
	userTest.UpdateProviderInformation("google", "102781991376288521134")
	if err != nil {
		panic(err)
	}

	userTest.UpdateLocation("Afurada")
	userTest.UpdateAllowsFriendRequests(false)
	userTest.UpdateAllowsRecommendations(false)

	userTest.RewardBadge(value.UserBadgeBrand)

	// Create user and get auto-generated ID
	_, err = userRepo.CreateUser(ctx, userTest)
	if err != nil {
		panic(err)
	}
	testID = userTest.ID

	return krayID, taikoID, testID
}

func BootstrapFriendships(ctx context.Context, friendshipRepo *repositories.FriendshipRepository, krayID, taikoID, testID int) {

	friendship := domain.NewFriendRequest(krayID, taikoID)
	friendship.Accept()

	err := friendshipRepo.CreateFriendship(ctx, friendship)
	if err != nil {
		panic(err)
	}

	friendship2 := domain.NewFriendRequest(taikoID, testID)
	friendship2.Accept()

	err = friendshipRepo.CreateFriendship(ctx, friendship2)
	if err != nil {
		panic(err)
	}
}

func BootstrapTranslations(ctx context.Context, translationRepo *repositories.DescriptionTranslationRepository) {

	translation, _ := domain.NewDescriptionTranslation(32901, "Há um ano, Sagiri Izumi tornou-se meio-irmã de Masamune Izumi. Mas a morte repentina de seus pais despedaça a nova família, fazendo com que Sagiri se isole do irmão e da sociedade.\n\nEnquanto cuida do que restou de sua família, Masamune ganha a vida como autor de light novels, com um pequeno problema: ele nunca conheceu sua aclamada ilustradora, Eromanga-sensei, famosa por desenhar as eróticas mais ousadas. Através de uma série de eventos embaraçosos, ele descobre que sua própria irmãzinha era sua parceira o tempo todo!\n\nÀ medida que novos personagens e desafios surgem, Masamune e Sagiri precisam enfrentar juntos a indústria de light novels. Eromanga-sensei acompanha o desenvolvimento do relacionamento deles e a luta para alcançar o sucesso; e, conforme Sagiri lentamente se liberta de sua timidez, por quanto tempo ela conseguirá esconder sua verdadeira personalidade do resto do mundo?\n\n[Escrito por MAL Rewrite]", 2)
	translation.Accept(2)

	err := translationRepo.CreateTranslation(ctx, translation)
	if err != nil {
		panic(err)
	}
}

func BootstrapReports(ctx context.Context, reportRepo *repositories.UserReportRepository) {

	report := domain.NewUserReport(value.ReportReasonIllegalActivities, 1, 2)
	err := reportRepo.CreateReport(ctx, report)
	if err != nil {
		panic(err)
	}
}

func BootstrapPostConversation(ctx context.Context, postRepo *repositories.PostRepository, krayID int, taikoID int, testID int) {

	post1 := domain.NewPost(strconv.Itoa(krayID), value.ParentTypeUser, "Bem vindos ao meu perfil ⭐!", krayID)

	createdPost1, err := postRepo.CreatePost(ctx, post1)
	if err != nil {
		panic(err)
	}

	post2 := domain.NewPost(createdPost1.ID, value.ParentTypePost, "Belo perfil, votos de bons dias! 👌", taikoID)

	createdPost2, err := postRepo.CreatePost(ctx, post2)
	if err != nil {
		panic(err)
	}
	if err := postRepo.AddReplyToPost(ctx, createdPost1.ID, createdPost2.ID); err != nil {
		panic(err)
	}

	post3 := domain.NewPost(createdPost1.ID, value.ParentTypePost, "Bem vindos ao website também! 😀", krayID)

	createdPost3, err := postRepo.CreatePost(ctx, post3)
	if err != nil {
		panic(err)
	}
	if err := postRepo.AddReplyToPost(ctx, createdPost1.ID, createdPost3.ID); err != nil {
		panic(err)
	}

	post4 := domain.NewPost(createdPost2.ID, value.ParentTypePost, "Bom dia para mim também !? :c", testID)

	createdPost4, err := postRepo.CreatePost(ctx, post4)
	if err != nil {
		panic(err)
	}
	if err := postRepo.AddReplyToPost(ctx, createdPost2.ID, createdPost4.ID); err != nil {
		panic(err)
	}
}

func BootstrapIndices(ctx context.Context, m *mongo.Database) {

	m.Collection("recommendations").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "receiver", Value: 1}, {Key: "seen", Value: 1}}},
		{Keys: bson.D{{Key: "initiator", Value: 1}, {Key: "receiver", Value: 1}, {Key: "anime", Value: 1}}},
	})
}

func BootstrapAnimeList(ctx context.Context, animeListRepo *repositories.AnimeListRepository, userID int) {
	testList := domain.NewPersonalAnimeList(userID)

	testEntry := domain.NewAnimeListItem(userID, 1, value.AnimeListItemStatusWatching)
	testEntry.UpdateProgress(5, 12)
	testEntry.UpdateNotes("Cowboy.")
	testEntry.AddRating(8, 9, 7)
	testList.AddListItem(*testEntry)

	if err := animeListRepo.SaveUserList(ctx, testList); err != nil {
		panic(err)
	}
}

func SuperFunnyListFiller3000(list *domain.UserAnimeList, userID, limit int) {
	for i := 1; i <= limit; i++ {
		entry := domain.NewAnimeListItem(userID, uint32(i), value.AnimeListItemStatusWatching)
		entry.UpdateProgress(uint32(i%12), 12)
		entry.UpdateNotes("Anime número " + strconv.Itoa(i))
		entry.AddRating(uint8((i%10)+1), uint8((i%10)+1), uint8((i%10)+1))
		list.AddListItem(*entry)
	}
}
