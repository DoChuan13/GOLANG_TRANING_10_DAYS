/*Package logger implements log logics.*/
package logger

import (
	"go.uber.org/zap"
)

var (
	sugaredLogger *zap.SugaredLogger
)

// AppLogVariables push Error Log related variables, file, objects,...
func AppLogVariables(messageID, level, variable, value string) {
	if sugaredLogger != nil {
		sugaredLogger.Infow(
			EmptyString, zap.String(LogTypeKey, AppLogType), zap.String(MessageIDKey, messageID),
			zap.String(LevelKey, level), zap.String(VariableNameKey, variable),
			zap.String(InstanceNameKey, InstanceName),
			zap.String(ValueNameKey, value),
		)
	}
}

// AppLogReadFile push Error Log number parameter error
func AppLogReadFile(messageID, level, path string) {
	if sugaredLogger != nil {
		sugaredLogger.Infow(
			EmptyString, zap.String(LogTypeKey, AppLogType), zap.String(MessageIDKey, messageID),
			zap.String(LevelKey, level), zap.String(InstanceNameKey, InstanceName),
			zap.String(FileNameKey, path),
		)
	}
}

// AppLogKafka push log about AWS S3
func AppLogKafka(messageID, level string) {
	if sugaredLogger != nil {
		sugaredLogger.Infow(
			EmptyString, zap.String(LogTypeKey, AppLogType), zap.String(MessageIDKey, messageID),
			zap.String(LevelKey, level), zap.String(InstanceNameKey, InstanceName),
		)
	}
}

// AppLog push application log
func AppLog(messageID, level string) {
	if sugaredLogger != nil {
		sugaredLogger.Infow(
			EmptyString, zap.String(LogTypeKey, AppLogType), zap.String(MessageIDKey, messageID),
			zap.String(LevelKey, level), zap.String(InstanceNameKey, InstanceName),
		)
	}
}

// AppLogDatabase push application log about database
func AppLogDatabase(messageID, level, endpoint, dbName string) {
	if sugaredLogger != nil {
		sugaredLogger.Infow(
			EmptyString, zap.String(LogTypeKey, AppLogType), zap.String(MessageIDKey, messageID),
			zap.String(LevelKey, level), zap.String(InstanceNameKey, InstanceName),
			zap.String(EndpointNameKey, endpoint),
			zap.String(DatabaseNameKey, dbName),
		)
	}
}

// AppLogQueryDatabase push application log about query database
func AppLogQueryDatabase(messageID, queryString string) {
	if sugaredLogger != nil {
		sugaredLogger.Infow(
			EmptyString, zap.String(LogTypeKey, AppLogType), zap.String(MessageIDKey, messageID),
			zap.String(LevelKey, ErrorLevelLog), zap.String(InstanceNameKey, InstanceName),
			zap.String(QueryKey, queryString),
		)
	}
}

// AppLogDatabaseStatusFalse log for database status is false
func AppLogDatabaseStatusFalse(messageID, kubun, hasshin, kei, endpoint, dbName string) {
	if sugaredLogger != nil {
		sugaredLogger.Infow(
			EmptyString, zap.String(LogTypeKey, AppLogType), zap.String(MessageIDKey, messageID),
			zap.String(LevelKey, ErrorLevelLog), zap.String(InstanceNameKey, InstanceName),
			zap.String(KubunKey, kubun), zap.String(HassinKey, hasshin), zap.String(KeiKey, kei),
			zap.String(EndpointNameKey, endpoint), zap.String(DatabaseNameKey, dbName),
		)
	}
}

// AppLogCannotGetTables log for cannot get tables from database
func AppLogCannotGetTables(messageID, kubun, hasshin, date, endpoint, dbName string) {
	if sugaredLogger != nil {
		sugaredLogger.Infow(
			EmptyString, zap.String(LogTypeKey, AppLogType), zap.String(MessageIDKey, messageID),
			zap.String(LevelKey, ErrorLevelLog), zap.String(InstanceNameKey, InstanceName),
			zap.String(KubunKey, kubun), zap.String(HassinKey, hasshin), zap.String(DateKey, date),
			zap.String(EndpointNameKey, endpoint), zap.String(DatabaseNameKey, dbName),
		)
	}
}

// AppLogCannotGetColumns log for cannot get columns
func AppLogCannotGetColumns(messageID, tableName, endpoint, dbName string) {
	if sugaredLogger != nil {
		sugaredLogger.Infow(
			EmptyString, zap.String(LogTypeKey, AppLogType), zap.String(MessageIDKey, messageID),
			zap.String(LevelKey, ErrorLevelLog), zap.String(InstanceNameKey, InstanceName),
			zap.String(TableNameKey, tableName), zap.String(EndpointNameKey, endpoint),
			zap.String(DatabaseNameKey, dbName),
		)
	}
}

// AppLogRecoverFail log for recover fail
func AppLogRecoverFail(messageID, tableName, sourceEndpoint, sourceDBName, targetEndpoint, targetDBName, step string) {
	if sugaredLogger != nil {
		sugaredLogger.Infow(
			EmptyString, zap.String(LogTypeKey, AppLogType), zap.String(MessageIDKey, messageID),
			zap.String(LevelKey, ErrorLevelLog), zap.String(InstanceNameKey, InstanceName),
			zap.String(TableNameKey, tableName), zap.String(SourceEndpointKey, sourceEndpoint),
			zap.String(SourceDatabaseNameKey, sourceDBName), zap.String(TargetEndpointKey, targetEndpoint),
			zap.String(TargetDatabaseNameKey, targetDBName), zap.String(StepKey, step),
		)
	}
}
