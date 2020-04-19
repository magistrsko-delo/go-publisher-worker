package Models

type NewMediaModel struct {
	MediaId int32
	Name string
	SiteName string
	Length int32
	Status int32
	Thumbnail string
	ProjectId int32
	AwsBucketWholeMedia string
	AwsStorageNameWholeMedia string
}
