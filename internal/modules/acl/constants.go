package acl

var AllResources []Resource = []Resource{
	GroupResource,
	UserResource,
}

var AllOperations []Operation = []Operation{
	ReadOps,
	WriteOps,
}
