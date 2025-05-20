package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xanzy/go-gitlab"

	"sonarqube-ouath-async/async"
	"sonarqube-ouath-async/flag"
	"sonarqube-ouath-async/jojogitlab"
	"sonarqube-ouath-async/log"
	"sonarqube-ouath-async/sonarqube"
)

func GitlabAsync2Sonar(c *gin.Context) {
	go async.ToSonar()
	c.JSON(200, gin.H{
		"data": "ok",
	})
}

func SyncAllProject(c *gin.Context) {
	data, err := async.SyncToSonar()
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(200, data)

}

func SyncProject(c *gin.Context) {
	var projectId = c.Param("projectId")
	data, err := async.SyncOneProjectToSonar(projectId)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(200, data)

}

func SyncUser(c *gin.Context) {
	username := c.Param("username")
	data, err := async.SyncOneUserToSonar(username)
	if err != nil {
		fmt.Println("err: ", err)
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(200, data)

}

func GetSonarProjectsFromPG(c *gin.Context) {
	projects, err := sonarqube.GetAllProjectsFromPG()
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(200, projects)
}

func GetSonarProjects(c *gin.Context) {
	projects, err := sonarqube.GetAllProjects()
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(200, projects)
}

func GetUsers(c *gin.Context) {
	users, err := sonarqube.GetAllUsers()
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(200, users)
}

func GetGitlabAllProjects(c *gin.Context) {
	var err error
	var gProjects []*gitlab.Project
	// Create GitLab client
	var gClient *gitlab.Client
	gClient, err = gitlab.NewClient(flag.Configuration.GitlabToken, gitlab.WithBaseURL(flag.Configuration.GitlabAddr))
	if err != nil {
		log.Logger.Errorf("Failed to create GitLab client: %v\n", err)
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	gProjects, err = jojogitlab.GetAllGitlabProjects(gClient)
	if err != nil {
		log.Logger.Errorf("Failed to get GitLab projects: %v\n", err)
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(200, gProjects)
}

func GetGitlabProjectMembers(c *gin.Context) {
	var projectId = c.Param("projectId")

	gClient, err := gitlab.NewClient(flag.Configuration.GitlabToken, gitlab.WithBaseURL(flag.Configuration.GitlabAddr))
	if err != nil {
		log.Logger.Errorf("Failed to create GitLab client: %v\n", err)
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	// Get gitlab project members
	gitlabMembers, err := jojogitlab.GetGitlabProjectMembers(gClient, projectId)
	if err != nil {
		log.Logger.Warningf("Failed to get gitlab project members: %v\n", err)
		log.Logger.Errorf("Failed to get GitLab project: %v\n", err)
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(200, gitlabMembers)
}
