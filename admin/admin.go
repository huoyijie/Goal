package admin

func (OperationLog) TableName() string {
	return "admin_operation_logs"
}

func (*OperationLog) GroupDynamicStrings() []string {
	return []string{"auth", "admin", "cdn"}
}

func (*OperationLog) ItemDynamicStrings() []string {
	return []string{"auth.role", "auth.session", "auth.user", "admin.operationlog", "cdn.resource"}
}
