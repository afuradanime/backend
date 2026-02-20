package repositories

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserReportRepository struct {
	collection        *mongo.Collection
	counterCollection *mongo.Collection
}

type ReportResult struct {
	Report   domain.UserReport
	Reporter *domain.User
	Target   *domain.User
}

func NewUserReportRepository(db *mongo.Database) *UserReportRepository {
	return &UserReportRepository{
		collection:        db.Collection("reports"),
		counterCollection: db.Collection("counters"),
	}
}

func (r *UserReportRepository) getNextSequence(ctx context.Context, name string) (int, error) {
	filter := bson.M{"_id": name}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var result Counter
	err := r.counterCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.Seq, nil
}

func (r *UserReportRepository) CreateReport(ctx context.Context, report *domain.UserReport) error {
	nextID, err := r.getNextSequence(ctx, "report_id")
	if err != nil {
		return err
	}
	report.ID = nextID
	_, err = r.collection.InsertOne(ctx, report)
	return err
}

func (r *UserReportRepository) GetReportByID(ctx context.Context, id int) (*domain.UserReport, error) {
	var report domain.UserReport
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *UserReportRepository) DeleteReport(ctx context.Context, id int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *UserReportRepository) GetReports(ctx context.Context, pageNumber, pageSize int) ([]ReportResult, utils.Pagination, error) {
	skip := (pageNumber - 1) * pageSize

	matchStage := bson.D{{Key: "$match", Value: bson.M{}}}

	lookupReporter := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "created_by",
		"foreignField": "_id",
		"as":           "reporter",
	}}}

	lookupTarget := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "target_user",
		"foreignField": "_id",
		"as":           "target",
	}}}

	// Count
	countCursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		bson.D{{Key: "$count", Value: "total"}},
	})
	if err != nil {
		return nil, utils.Pagination{}, err
	}
	var countResult []struct {
		Total int `bson:"total"`
	}
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, utils.Pagination{}, err
	}
	total := 0
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	pipeline := mongo.Pipeline{
		matchStage,
		lookupReporter,
		lookupTarget,
		bson.D{{Key: "$sort", Value: bson.M{"_id": -1}}},
		bson.D{{Key: "$skip", Value: int64(skip)}},
		bson.D{{Key: "$limit", Value: int64(pageSize)}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var results []struct {
		domain.UserReport `bson:",inline"`
		Reporter          []domain.User `bson:"reporter"`
		Target            []domain.User `bson:"target"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, utils.Pagination{}, err
	}

	reports := make([]ReportResult, len(results))
	for i, r := range results {
		reports[i] = ReportResult{Report: r.UserReport}
		if len(r.Reporter) > 0 {
			reports[i].Reporter = &r.Reporter[0]
		}
		if len(r.Target) > 0 {
			reports[i].Target = &r.Target[0]
		}
	}

	totalPages := (total + pageSize - 1) / pageSize
	return reports, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *UserReportRepository) GetReportsByTarget(ctx context.Context, targetUserID int, pageNumber, pageSize int) ([]domain.UserReport, utils.Pagination, error) {
	skip := (pageNumber - 1) * pageSize
	filter := bson.M{"target_user": targetUserID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"_id": -1}),
	)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var reports []domain.UserReport
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, utils.Pagination{}, err
	}

	totalPages := (int(total) + pageSize - 1) / pageSize
	return reports, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *UserReportRepository) HasReported(ctx context.Context, reporterID, targetUserID int) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"created_by":  reporterID,
		"target_user": targetUserID,
	})
	return count > 0, err
}
