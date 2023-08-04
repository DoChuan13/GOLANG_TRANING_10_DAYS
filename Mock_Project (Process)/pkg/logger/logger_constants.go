package logger

const (
	// ServiceName service name
	ServiceName = "new-recover-data"
	// LogAppPort port for log system
	LogAppPort = 9999
	// LogAppIP IP for log system
	LogAppIP = "0.0.0.0"
)

const (
	// EmptyString empty character
	EmptyString = ""
)

// InstanceName instance name
var InstanceName string

const (
	// InfoAppNewRecoverDataFinished log ID for app new recover data finished
	InfoAppNewRecoverDataFinished = "INF909"
	// InfoAppStartNewRecoverData log ID for app start new recover data
	InfoAppStartNewRecoverData = "INF908"
)

const (
	// ErrorCodeIncorrectArgumentNumber error code for case number parameters incorrect
	ErrorCodeIncorrectArgumentNumber = "ERR800"
	// ErrorCodeVariableNotFound error code for case variable not found
	ErrorCodeVariableNotFound = "ERR801"
	// ErrorCodeInvalidVariable error code for invalid variable
	ErrorCodeInvalidVariable = "ERR802"
	// ErrorCodeWrongDataFormat error code for wrong data format
	ErrorCodeWrongDataFormat = "ERR803"
	// ErrorCodeReadFileConfigFail error code for case error read file environment variables
	ErrorCodeReadFileConfigFail = "ERR804"
	// ErrorCodeParseJSONFail error code for case error parse json fail
	ErrorCodeParseJSONFail = "ERR805"
	// ErrorCodeInitS3ConnectionFail error code for case error cannot connect to AWS S3
	ErrorCodeInitS3ConnectionFail = "ERR806"
	// ErrorCodeInitDatabaseConnectionFail error code for case error cannot connect to Database
	ErrorCodeInitDatabaseConnectionFail = "ERR808"
	// ErrorCodeExecuteQueryFail error code for case  error execute query string fail
	ErrorCodeExecuteQueryFail = "ERR810"
	// ErrorCodeDatabaseStatusIsFalse error code of case database status is false
	ErrorCodeDatabaseStatusIsFalse = "ERR886"
	// ErrorCodeCannotGetTablesFromDatabase error code of case cannot get tables
	ErrorCodeCannotGetTablesFromDatabase = "ERR887"
	// ErrorCodeCannotGetColumns error code of case cannot get columns
	ErrorCodeCannotGetColumns = "ERR888"
	// ErrorCodeRecoverFail error code of case recover fail
	ErrorCodeRecoverFail = "ERR889"
)

const (
	// AppLogType APP LOG type
	AppLogType = "TCK_APPLOG"
	// InfoLevelLog level log info
	InfoLevelLog = "INFO"
	// ErrorLevelLog level log error
	ErrorLevelLog = "ERROR"
	// LogTypeKey log type key
	LogTypeKey = "LogType"
	// MessageIDKey message ID key
	MessageIDKey = "MessageID"
	// LevelKey level key
	LevelKey = "Level"
	// ServiceNameKey service name key
	ServiceNameKey = "ServiceName"
	// InstanceNameKey instance name key
	InstanceNameKey = "InstanceName"
	// VariableNameKey variable name key
	VariableNameKey = "VariableName"
	// FileNameKey file name key
	FileNameKey = "FileName"
	// ValueNameKey value name key
	ValueNameKey = "Value"
	// EndpointNameKey endpoint name key
	EndpointNameKey = "EndpointName"
	// DatabaseNameKey database name key
	DatabaseNameKey = "DatabaseName"
	// QueryKey query statement key
	QueryKey = "Sql"
	// RegionKey region key
	RegionKey = "S3Region"
	// BucketKey bucket key
	BucketKey = "S3Bucket"
	// EC2HostnameKey EC2 instance name key
	EC2HostnameKey = "EC2InstanceHostname"
	// SocketHostKey socket host key
	SocketHostKey = "SocketHost"
	// SocketPortKey socket port key
	SocketPortKey = "SocketPort"
	// CommandNameKey command name key
	CommandNameKey = "CommandName"
	// QueueMessageKey queue message key
	QueueMessageKey = "QueueMessage"
	// QuoteCodeKey quote code key
	QuoteCodeKey = "QuoteCode"
	// OldQuoteCodeKey old quote code key
	OldQuoteCodeKey = "OldQuoteCode"
	// NewQuoteCodeKey new quote code key
	NewQuoteCodeKey = "NewQuoteCode"
	// FromDateKey from date key
	FromDateKey = "FromDate"
	// ToDateKey to date key
	ToDateKey = "ToDate"
	// KubunKey kubun key
	KubunKey = "Kubun"
	// HassinKey hassin key
	HassinKey = "Hasshin"
	// KeiKey kei key
	KeiKey = "Kei"
	// DateKey date key
	DateKey = "Date"
	// TableNameKey table name key
	TableNameKey = "TableName"
	// SourceEndpointKey source endpoint key
	SourceEndpointKey = "SourceEndpoint"
	// SourceDatabaseNameKey source database name key
	SourceDatabaseNameKey = "SourceDatabaseName"
	// TargetEndpointKey target endpoint key
	TargetEndpointKey = "TargetEndpoint"
	// TargetDatabaseNameKey target database name key
	TargetDatabaseNameKey = "TargetDatabaseName"
	// StepKey step key
	StepKey = "Step"
)
