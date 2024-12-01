package main

import "time"

type Log struct {
	FileName                    string  `json:"file_name" bson:"file_name"`
	FileUpload                  LogData `json:"file_upload" bson:"file_upload"`
	S3TriggerCompleted          LogData `json:"s3_trigger_completed" bson:"s3_trigger_completed"`
	ParsedToJson                LogData `json:"parsed_to_json" bson:"parsed_to_json"`
	SourceSchemaMetadata        LogData `json:"source_schema_metadata" bson:"source_schema_metadata"`
	TargetSchemaMetadata        LogData `json:"target_schema_metadata" bson:"target_schema_metadata"`
	AssociationMappingGenerated LogData `json:"association_mapping_generated" bson:"association_mapping_generated"`
	DataInjectionCompleted      LogData `json:"data_injection_completed" bson:"data_injection_completed"`
}

type LogData struct {
	Status    bool      `json:"status" bson:"status"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
