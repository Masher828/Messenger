package mongo_common_repo

import (
	"context"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func GetSingleDocumentByFilter(log *zap.SugaredLogger, collectionName string, filter map[string]interface{}, data interface{}) error {
	db := system.MessengerContext.MongoDB

	collection := db.Database(system.MongoDatabaseName).Collection(collectionName)

	err := collection.FindOne(context.TODO(), filter).Decode(data)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func InsertDocument(log *zap.SugaredLogger, collectionName string, data interface{}) error {
	db := system.MessengerContext.MongoDB

	collection := db.Database(system.MongoDatabaseName).Collection(collectionName)

	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func UpdateDocumentById(log *zap.SugaredLogger, collectionName, documentId string, dataToUpdate map[string]interface{}) error {
	db := system.MessengerContext.MongoDB

	collection := db.Database(system.MongoDatabaseName).Collection(collectionName)

	dataToBeUpdated := map[string]interface{}{"$set": dataToUpdate}

	_, err := collection.UpdateByID(context.TODO(), documentId, dataToBeUpdated)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func GetDocumentCountsByFilter(log *zap.SugaredLogger, collectionName string, filter map[string]interface{}) (int64, error) {
	db := system.MessengerContext.MongoDB

	collection := db.Database(system.MongoDatabaseName).Collection(collectionName)

	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Errorln(err)
		return 0, err
	}

	return count, nil

}

func GetDocumentsWithFilter(log *zap.SugaredLogger, collectionName string, filter map[string]interface{}, offset, limit int64, data interface{}) error {

	db := system.MessengerContext.MongoDB

	collection := db.Database(system.MongoDatabaseName).Collection(collectionName)

	opts := options.FindOptions{}
	if offset > 0 {
		opts.SetSkip(offset)
	}

	if limit > 0 {
		opts.SetLimit(limit)
	}

	opts.SetSort(map[string]interface{}{"createdOn": -1})

	cursor, err := collection.Find(context.TODO(), filter, &opts)
	if err != nil {
		log.Errorln(err)
		return err
	}

	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &data)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return err
}

func GetSelectedFieldsDocumentsWithFilter(log *zap.SugaredLogger, collectionName string, selectedFields, filter map[string]interface{}, offset, limit int64, data interface{}) error {

	db := system.MessengerContext.MongoDB

	collection := db.Database(system.MongoDatabaseName).Collection(collectionName)

	opts := options.FindOptions{}
	if offset > 0 {
		opts.SetSkip(offset)
	}

	if limit > 0 {
		opts.SetLimit(limit)
	}

	opts.SetSort(map[string]interface{}{"createdOn": -1})
	opts.SetProjection(selectedFields)

	cursor, err := collection.Find(context.TODO(), filter, &opts)
	if err != nil {
		log.Errorln(err)
		return err
	}

	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), data)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return err
}

func UpsertDocumentCustomQuery(log *zap.SugaredLogger, collectionName string, filter, updateQuery map[string]interface{}) error {

	db := system.MessengerContext.MongoDB

	collection := db.Database(system.MongoDatabaseName).Collection(collectionName)

	opts := options.UpdateOptions{}
	opts.SetUpsert(true)

	_, err := collection.UpdateMany(context.TODO(), filter, updateQuery, &opts)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func DeleteSingleDocumentByFilter(log *zap.SugaredLogger, collectionName string, filter map[string]interface{}) error {
	db := system.MessengerContext.MongoDB

	collection := db.Database(system.MongoDatabaseName).Collection(collectionName)

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}
