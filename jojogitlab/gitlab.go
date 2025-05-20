package jojogitlab

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"sonarqube-ouath-async/flag"
	"sonarqube-ouath-async/log"

	"github.com/go-resty/resty/v2"
	"github.com/xanzy/go-gitlab"
)

// Returns group names, already separated by '.'
func GetAllGroups() (groups map[string]int) {
	// TODO
	type GroupInfo struct {
		ID       int    `json:"id"`
		FullPath string `json:"full_path"`
	}
	ret, err := requestGit(flag.GitlabAddr + "/api/v4/groups")
	if err != nil {
		return
	}
	var res = make(map[string]int)
	for _, r := range ret {
		var tmp []GroupInfo
		err = json.Unmarshal(r, &tmp)
		if err != nil {
			return
		}
		for _, g := range tmp {
			res[strings.ReplaceAll(g.FullPath, "/", ".")] = g.ID
		}
	}

	return res
}

// Returns list of usernames
func GetGroupMembers(gid int) (response map[int]string) {
	// TODO
	type UserInfo struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
	}
	result := make(map[int]string)

	ret, err := requestGit(flag.GitlabAddr + "/api/v4/groups/" + strconv.Itoa(gid) + "/members")
	if err != nil {
		return
	}
	for _, r := range ret {
		var tmp []UserInfo
		err = json.Unmarshal(r, &tmp)
		if err != nil {
			return
		}
		for _, u := range tmp {
			result[u.ID] = u.Username
		}
	}
	return result
}

// Returns all user IDs and emails
func GetAllUsers() (response map[int]string) {
	// TODO
	type UserInfo struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	result := make(map[int]string)

	ret, err := requestGit(flag.GitlabAddr + "/api/v4/users")
	if err != nil {
		return
	}
	for _, r := range ret {
		var tmp []UserInfo
		err = json.Unmarshal(r, &tmp)
		if err != nil {
			return
		}
		for _, u := range tmp {
			result[u.ID] = u.Email
		}
	}
	return result
}

func requestGit(url string) (ret [][]byte, err error) {
	var page = 1
	client := resty.New()
	for {
		resp, err := client.R().
			SetQueryParams(map[string]string{
				"page":     strconv.Itoa(page),
				"per_page": "100",
			}).
			SetHeader("Accept", "application/json").
			SetHeader("Authorization", "Bearer "+flag.GitlabToken).
			Get(url)

		if err != nil {
			log.Logger.Error("Gitlab request error")
			return ret, err
		}
		// Return empty array [] to break the loop
		if len(resp.Body()) <= 2 {
			break
		}
		ret = append(ret, resp.Body())

		page++
	}

	return
}

func GetAllProjects() (projects []*gitlab.Project, err error) {
	// Create GitLab client
	client, err := gitlab.NewClient(flag.Configuration.GitlabToken, gitlab.WithBaseURL(flag.Configuration.GitlabAddr))
	if err != nil {
		log.Logger.Error("Failed to create GitLab client: %v\n", err)
		return
	}
	projects, err = GetAllGitlabProjects(client)
	if err != nil {
		return
	}
	return
}

func GetAllGitlabProjects(gitlabClient *gitlab.Client) ([]*gitlab.Project, error) {
	var allProjects []*gitlab.Project
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100, Page: 1},
	}

	for {
		projects, resp, err := gitlabClient.Projects.ListProjects(opt)
		if err != nil {
			return nil, err
		}
		allProjects = append(allProjects, projects...)

		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opt.Page = resp.NextPage

		time.Sleep(1 * time.Second)
	}

	return allProjects, nil
}

func GetGitlabProjectInfo(gitlabClient *gitlab.Client, projectId string) (*gitlab.Project, error) {
	project, response, err := gitlabClient.Projects.GetProject(projectId, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		body, _ := io.ReadAll(response.Body)
		err = errors.New(fmt.Sprintf("response.StatusCode: %v, body: %v", response.StatusCode, string(body)))
		return nil, err
	}

	return project, nil
}

func GetAllGitlabUsers(gitlabClient *gitlab.Client) ([]*gitlab.User, error) {
	var allUsers []*gitlab.User
	opt := &gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100, Page: 1},
	}

	for {
		users, resp, err := gitlabClient.Users.ListUsers(opt)
		if err != nil {
			return nil, err
		}
		allUsers = append(allUsers, users...)

		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opt.Page = resp.NextPage
		time.Sleep(1 * time.Second)
	}

	return allUsers, nil
}

func GetUseridByUsername(gitlabClient *gitlab.Client, username string) (*gitlab.User, error) {
	opt := &gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100, Page: 1},
		Username:    &username,
	}

	users, _, err := gitlabClient.Users.ListUsers(opt)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		err = errors.New(fmt.Sprintf("No user found with username: %s", username))
		return nil, err
	}

	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("No user found with username: %s", username))
}

func GetGitlabProjectMembers(gitlabClient *gitlab.Client, projectID interface{}) (map[string]*gitlab.ProjectMember, error) {
	opt := &gitlab.ListProjectMembersOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100, Page: 1},
	}

	var allMembers []*gitlab.ProjectMember
	for {
		members, resp, err := gitlabClient.ProjectMembers.ListAllProjectMembers(projectID, opt)
		//members, resp, err := gitlabClient.ProjectMembers.ListProjectMembers(projectID, opt)
		if err != nil {
			return nil, err
		}
		allMembers = append(allMembers, members...)

		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opt.Page = resp.NextPage

		time.Sleep(1 * time.Second)
	}

	var data = make(map[string]*gitlab.ProjectMember)
	for _, member := range allMembers {
		data[member.Username] = member
	}

	return data, nil
}

func ChangeGitlabProjectsToMap(projects []*gitlab.Project) (mapGitlabProjects map[string]*gitlab.Project) {
	mapGitlabProjects = make(map[string]*gitlab.Project)
	for _, p := range projects {
		var key = strings.ReplaceAll(p.Namespace.FullPath, "/", ".") + fmt.Sprintf(":%s", p.Path) + fmt.Sprintf(":%v", p.ID)
		//fmt.Println("gitlab.Project key:", key)
		mapGitlabProjects[key] = p
	}
	return
}
