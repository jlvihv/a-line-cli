package model

/*
Artifactory 构建物
*/
type Artifactory struct {
	Name string
	Url  string
}

/*
Report 构建物报告
*/
type Report struct {
	Id   int
	Url  string
	Type int
}

type ActionResult struct {
	Artifactorys []Artifactory
	Reports      []Report
}
