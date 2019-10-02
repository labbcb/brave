package mongo

import (
	"github.com/labbcb/brave/search"
	"github.com/labbcb/brave/variant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB extends Mongo client to address business model (genomic variants)
type DB struct {
	client   *mongo.Client
	database string
}

// Connect creates connection with MongoDB
func Connect(uri string, database string) (*DB, error) {
	c, err := mongo.Connect(nil, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &DB{client: c, database: database}, nil
}

// Save stores variant to Mongo database
func (db *DB) Save(v *variant.Variant) error {
	_, err := db.client.Database(db.database).Collection("variants").InsertOne(nil, v)
	if err != nil {
		return err
	}
	return nil
}

// Search is the main method to search for variants.
func (db *DB) Search(i *search.Input) (*search.Response, error) {
	var filters bson.A
	for _, q := range i.Queries {
		var fq bson.A
		if q.AssemblyID != "" {
			fq = append(fq, bson.D{{"assemblyId", q.AssemblyID}})
		}
		if q.GeneSymbol != "" {
			fq = append(fq, bson.D{{"geneSymbol", bson.D{{"$all", bson.A{q.GeneSymbol}}}}})
		}
		if q.DatasetID != "" {
			fq = append(fq, bson.D{{"datasetId", q.DatasetID}})
		}
		if q.SnpID != "" {
			fq = append(fq, bson.D{{"snpIds", bson.D{{"$all", bson.A{q.SnpID}}}}})
		}
		if q.ReferenceName != "" && q.Start != 0 && q.End != 0 {
			fq = append(fq, bson.D{{"$and", bson.A{
				bson.D{{"referenceName", q.ReferenceName}},
				bson.D{{"start", bson.D{{"$gte", q.Start}}}},
				bson.D{{"start", bson.D{{"$lte", q.End}}}},
			}}})
		} else if q.ReferenceName != "" && q.Start != 0 {
			fq = append(fq, bson.D{{"$and", bson.A{
				bson.D{{"referenceName", q.ReferenceName}},
				bson.D{{"start", q.Start}},
			}}})
		}
		filters = append(filters, bson.D{{"$and", fq}})
	}

	filter := bson.D{}
	if len(filters) > 0 {
		filter = bson.D{{"$or", filters}}
	}

	cur, err := db.client.Database(db.database).Collection("variants").
		Find(nil, filter, &options.FindOptions{Limit: &i.Length, Skip: &i.Start})
	if err != nil {
		return nil, err
	}

	var variants []*variant.Variant
	if err := cur.All(nil, &variants); err != nil {
		return nil, err
	}

	if variants == nil {
		variants = []*variant.Variant{}
	}

	total, err := db.client.Database(db.database).Collection("variants").CountDocuments(nil, bson.D{})
	if err != nil {
		return nil, err
	}

	filtered, err := db.client.Database(db.database).Collection("variants").CountDocuments(nil, filter)
	if err != nil {
		return nil, err
	}

	return &search.Response{Draw: i.Draw, Variants: variants, RecordsTotal: total, RecordsFiltered: filtered}, nil
}

// Remove removes variants from database given a dataset ID and/or assembly ID.
// If both are zero value them it deletes all variants.
func (db *DB) Remove(datasetID string, assemblyID string) error {
	var filters bson.A
	if datasetID != "" {
		filters = append(filters, bson.D{{"datasetId", datasetID}})
	}
	if assemblyID != "" {
		filters = append(filters, bson.D{{"assemblyId", assemblyID}})
	}

	filter := bson.D{}
	if len(filters) > 0 {
		filter = bson.D{{"$and", filters}}
	}

	_, err := db.client.Database(db.database).Collection("variants").DeleteMany(nil, filter)
	if err != nil {
		return err
	}

	return nil
}
