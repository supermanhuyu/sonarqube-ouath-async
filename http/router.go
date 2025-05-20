package http

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(v1 *gin.RouterGroup) {
	v1.POST("/async", GitlabAsync2Sonar)

	// Sync operations, update sonar data
	sync := v1.Group("/sync")
	{
		sync.POST("/all", SyncAllProject)
		sync.POST("/project/:projectId", SyncProject)
		sync.POST("/user/:username", SyncUser)
	}

	// Query sonar related data
	sonar := v1.Group("/sonar")
	{
		//sonar.GET("/projects", GetSonarProjects)
		sonar.GET("/projects/list", GetSonarProjectsFromPG)
		sonar.GET("/users/list", GetUsers)
	}

	// Query gitlab related data
	gitlab := v1.Group("/gitlab")
	{
		gitlab.GET("/projects/list", GetGitlabAllProjects)
		gitlab.GET("/users/:projectId", GetGitlabProjectMembers)
	}
}
