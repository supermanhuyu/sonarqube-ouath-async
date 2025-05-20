package async

import (
	"fmt"
	"strings"
	"time"

	"sonarqube-ouath-async/flag"
	"sonarqube-ouath-async/jojogitlab"
	"sonarqube-ouath-async/log"
	"sonarqube-ouath-async/sonarqube"

	"github.com/xanzy/go-gitlab"
)

func ToSonar() {
	gitlabGroups := jojogitlab.GetAllGroups()
	gitlabIdEmailMap := jojogitlab.GetAllUsers()
	groupNameIdMap := sonarqube.SearchGroup()
	emailLoginMap := sonarqube.SearchUser()
	loginEmailMap := make(map[string]string)

	for email, login := range emailLoginMap {
		loginEmailMap[login] = email
	}

	for gitlabGroup, gid := range gitlabGroups {
		// Replace '/' with '.' in group names if needed
		_, ok := groupNameIdMap[gitlabGroup]
		if !ok {
			sonarqube.CreatGroup(gitlabGroup, "")
		}
		// Get users of the group
		gitlabIdUserNameMap := jojogitlab.GetGroupMembers(gid)

		gitlabEmailIdMap := make(map[string]int)
		for id, _ := range gitlabIdUserNameMap {
			_, existGitUser := gitlabIdEmailMap[id]
			if existGitUser {
				email := gitlabIdEmailMap[id]
				gitlabEmailIdMap[email] = id
			}
		}

		sonarLoginNameMap := sonarqube.GetGroupUsers(gitlabGroup)
		sonarEmailLoginMap := make(map[string]string)
		sonarLoginEmailMap := make(map[string]string)
		for login, _ := range sonarLoginNameMap {
			_, existSonarUser := loginEmailMap[login]
			if existSonarUser {
				email := loginEmailMap[login]
				sonarEmailLoginMap[email] = login
				sonarLoginEmailMap[login] = email
			}
		}

		for email, _ := range gitlabEmailIdMap {
			_, unAddUser := sonarEmailLoginMap[email]
			_, existSonaruser := emailLoginMap[email]
			if !unAddUser && existSonaruser {
				// If the user exists and is not currently included
				login := emailLoginMap[email]
				sonarqube.AddUserToGroup(gitlabGroup, login)
				fmt.Println("AddUserToGroup " + gitlabGroup + " " + login)
			}
		}

		for sonarEmail, login := range sonarEmailLoginMap {
			_, ok := gitlabEmailIdMap[sonarEmail]
			if !ok {
				// Remove users who have lost permissions
				sonarqube.RemoveUserToGroup(gitlabGroup, login)
				fmt.Println("RemoveUserToGroup " + gitlabGroup + " " + login)
			}
		}
	}

	// Query user groups again
	groupNameIdMap = sonarqube.SearchGroup()
	// Query permission template list
	namespaceIdMap := sonarqube.SearchTemplates()

	for groupName, _ := range groupNameIdMap {
		_, existTemplate := namespaceIdMap[groupName]
		var templateId string
		if !existTemplate {
			// If the template does not exist, create it with a suffix (:.*)
			templateId = sonarqube.CreateTemplate(groupName, "", groupName+":.*")
		} else {
			templateId = namespaceIdMap[groupName]
		}

		subGroupName := GetAllGroupName(groupName)
		for _, group := range subGroupName {
			sonarqube.AddGroupToTemplate(templateId, group, "user")
			sonarqube.AddGroupToTemplate(templateId, group, "codeviewer")
			sonarqube.AddGroupToTemplate(templateId, group, "issueadmin")
			sonarqube.AddGroupToTemplate(templateId, group, "securityhotspotadmin")
			sonarqube.AddGroupToTemplate(templateId, group, "admin")
		}
		sonarqube.AddGroupToTemplate(templateId, "sonar-administrators", "admin")
	}
}

func RunSyncToSonar() {
	_, err := SyncToSonar()
	if err != nil {
		fmt.Println("SyncToSonar error:", err)
		return
	}
}
func SyncToSonar() (data []PermissionData, err error) {
	fmt.Println("start sync gitlab project member permission to sonar project ...")

	var gProjects []*gitlab.Project
	//var gUsers []*gitlab.User
	// Create GitLab client
	var gClient *gitlab.Client
	gClient, err = gitlab.NewClient(flag.Configuration.GitlabToken, gitlab.WithBaseURL(flag.Configuration.GitlabAddr))
	if err != nil {
		log.Logger.Errorf("Failed to create GitLab client: %v\n", err)
		return
	}
	gProjects, err = jojogitlab.GetAllGitlabProjects(gClient)
	if err != nil {
		log.Logger.Errorf("Failed to get GitLab projects: %v\n", err)
		return
	}
	mapGitlabProjects := jojogitlab.ChangeGitlabProjectsToMap(gProjects)
	//fmt.Println("mapGitlabProjects: ", mapGitlabProjects)

	sonarProjects, err := sonarqube.GetAllProjectsFromPG()
	if err != nil {
		log.Logger.Errorf("Failed to get Sonar projects: %v\n", err)
		return
	}

	sonarUsers, err := sonarqube.GetAllUsers()
	if err != nil {
		log.Logger.Errorf("Failed to get Sonar users: %v\n", err)
		return
	}
	var mapSonarUsers = make(map[string]sonarqube.User)
	var mapLoginSonarUsers = make(map[string]sonarqube.User)
	for _, s := range sonarUsers {
		//fmt.Println("s.ExternalIdentity:", s.ExternalIdentity)
		mapSonarUsers[s.ExternalIdentity] = s
		mapLoginSonarUsers[s.Login] = s

	}

	for _, sProject := range sonarProjects {
		fmt.Println("\nsProject.Kee:", sProject.Kee)
		gitlabProject, exists := mapGitlabProjects[sProject.Kee]
		//fmt.Println("exists:", exists)
		// GitLab project exists in Sonar project
		if exists {
			// Get members of gitlab project
			gitlabMembers, err := jojogitlab.GetGitlabProjectMembers(gClient, gitlabProject.ID)
			if err != nil {
				log.Logger.Warningf("Failed to get gitlab project members: %v\n", err)
				continue
			}

			// Grant permissions to each member in sonar project
			for _, gitlabUser := range gitlabMembers {
				fmt.Println("gitlabUser.Username:", gitlabUser.Username)
				if gitlabUser.Username == "" || gitlabUser.Username == "root" || gitlabUser.Username == "admin" {
					continue
				}
				sonarUser, e := mapSonarUsers[gitlabUser.Username]
				if e {
					// Grant permissions to users existing in both gitlab and sonar
					addData, err := AddUserPermissionToProject(sProject.Kee, sonarUser.Login, int(gitlabUser.AccessLevel))
					if err != nil {
						log.Logger.Warningf("Failed to grant gitlab project member permission: %v\n", err)
						continue
					}
					data = append(data, addData...)
				}
			}

			projectUsers, _err := sonarqube.GetProjectUsers(sProject.Kee)
			if _err != nil {
				log.Logger.Warningf("Failed to get sonar project users: %v\n", err)
				continue
			}
			// Remove permissions for users in sonar project but not in gitlab project
			for _, pUser := range projectUsers.Users {
				// If Permissions is empty, all following users have no permissions, can exit.
				if len(pUser.Permissions) == 0 {
					continue
				}
				sonarUser, exist := mapLoginSonarUsers[pUser.Login]
				if !exist {
					continue
				}
				_, e := gitlabMembers[sonarUser.ExternalIdentity]
				if !e {
					//addData, err := RemoveUserAllPermission(sProject.Kee, pUser.Login, Print)
					addData, err := RemoveUserCertainPermission(sProject.Kee, sonarUser.Login, pUser.Permissions)
					if err != nil {
						log.Logger.Warningf("Failed to remove sonar project member permissions: %v\n", err)
						continue
					}
					data = append(data, addData...)
				}
			}
		}

		addData, err := UpdateGroupPermissionToProject(sProject.Kee)
		if err != nil {
			log.Logger.Warningf("Failed to update sonar project group permissions: %v\n", err)
		}
		data = append(data, addData...)
	}
	return
}

func SyncOneProjectToSonar(projectId string) (data []PermissionData, err error) {
	fmt.Println("start sync gitlab project member permission to sonar project ...")

	// Create GitLab client
	var gClient *gitlab.Client
	gClient, err = gitlab.NewClient(flag.Configuration.GitlabToken, gitlab.WithBaseURL(flag.Configuration.GitlabAddr))
	if err != nil {
		log.Logger.Errorf("Failed to create GitLab client: %v\n", err)
		return
	}

	// Get gitlab project
	gitlabProject, err := jojogitlab.GetGitlabProjectInfo(gClient, projectId)
	if err != nil {
		log.Logger.Warningf("Failed to get gitlab project: %v\n", err)
		return
	}
	// sonar kee
	var sProjectKee = strings.ReplaceAll(gitlabProject.Namespace.FullPath, "/", ".") + fmt.Sprintf(":%s", gitlabProject.Path) + fmt.Sprintf(":%v", gitlabProject.ID)
	fmt.Println("sonar ProjectKee:", sProjectKee)

	// Get members of gitlab project
	gitlabMembers, err := jojogitlab.GetGitlabProjectMembers(gClient, projectId)
	if err != nil {
		log.Logger.Warningf("Failed to get gitlab project members: %v\n", err)
		return
	}

	sonarUsers, err := sonarqube.GetAllUsers()
	if err != nil {
		log.Logger.Errorf("Failed to get Sonar users: %v\n", err)
		return
	}
	var mapSonarUsers = make(map[string]sonarqube.User)
	var mapLoginSonarUsers = make(map[string]sonarqube.User)
	for _, s := range sonarUsers {
		//fmt.Println("s.ExternalIdentity:", s.ExternalIdentity)
		mapSonarUsers[s.ExternalIdentity] = s
		mapLoginSonarUsers[s.Login] = s
	}

	// GitLab project exists in Sonar project

	// Grant permissions to each member in sonar project
	for _, gitlabUser := range gitlabMembers {
		fmt.Println("gitlabUser.Username:", gitlabUser.Username)
		if gitlabUser.Username == "" || gitlabUser.Username == "root" || gitlabUser.Username == "admin" {
			continue
		}
		sonarUser, e := mapSonarUsers[gitlabUser.Username]
		if e {
			// Grant permissions to users existing in both gitlab and sonar
			addData, err := AddUserPermissionToProject(sProjectKee, sonarUser.Login, int(gitlabUser.AccessLevel))
			if err != nil {
				log.Logger.Warningf("Failed to grant gitlab project member permission: %v\n", err)
				continue
			}
			data = append(data, addData...)
		}
	}

	projectUsers, _err := sonarqube.GetProjectUsers(sProjectKee)
	if _err != nil {
		log.Logger.Warningf("Failed to get sonar project users: %v\n", err)
		return
	}
	// Remove permissions for users in sonar project but not in gitlab project
	for _, pUser := range projectUsers.Users {
		// If Permissions is empty, all following users have no permissions, can exit.
		if len(pUser.Permissions) == 0 {
			continue
		}
		sonarUser, exist := mapLoginSonarUsers[pUser.Login]
		if !exist {
			continue
		}
		_, e := gitlabMembers[sonarUser.ExternalIdentity]
		if !e {
			addData, err := RemoveUserCertainPermission(sProjectKee, sonarUser.Login, pUser.Permissions)
			if err != nil {
				log.Logger.Warningf("Failed to remove sonar project member permissions: %v\n", err)
				continue
			}
			data = append(data, addData...)
		}
	}

	addData, err := UpdateGroupPermissionToProject(sProjectKee)
	if err != nil {
		log.Logger.Warningf("Failed to update sonar project group permissions: %v\n", err)
		return
	}
	data = append(data, addData...)

	return
}

func SyncOneUserToSonar(username string) (data []PermissionData, err error) {
	fmt.Printf("start sync gitlab project member [%s] permission to sonar project ...\n", username)

	var gProjects []*gitlab.Project
	// Create GitLab client
	var gClient *gitlab.Client
	gClient, err = gitlab.NewClient(flag.Configuration.GitlabToken, gitlab.WithBaseURL(flag.Configuration.GitlabAddr))
	if err != nil {
		log.Logger.Errorf("Failed to create GitLab client: %v\n", err)
		return
	}
	gProjects, err = jojogitlab.GetAllGitlabProjects(gClient)
	if err != nil {
		log.Logger.Errorf("Failed to get GitLab projects: %v\n", err)
		return
	}
	mapGitlabProjects := jojogitlab.ChangeGitlabProjectsToMap(gProjects)
	//fmt.Println("mapGitlabProjects: ", mapGitlabProjects)

	sonarProjects, err := sonarqube.GetAllProjectsFromPG()
	if err != nil {
		log.Logger.Errorf("Failed to get Sonar projects: %v\n", err)
		return
	}

	sonarUsers, err := sonarqube.GetAllUsers()
	if err != nil {
		log.Logger.Errorf("Failed to get Sonar users: %v\n", err)
		return
	}
	var mapSonarUsers = make(map[string]sonarqube.User)
	for _, s := range sonarUsers {
		//fmt.Println("s.ExternalIdentity:", s.ExternalIdentity)
		mapSonarUsers[s.ExternalIdentity] = s
	}

	for _, sProject := range sonarProjects {
		fmt.Println("\nsProject.Kee:", sProject.Kee)
		gitlabProject, exists := mapGitlabProjects[sProject.Kee]
		//fmt.Println("exists:", exists)
		// GitLab project exists in Sonar project
		if exists {
			// Get members of gitlab project
			gitlabMembers, err := jojogitlab.GetGitlabProjectMembers(gClient, gitlabProject.ID)
			if err != nil {
				log.Logger.Warningf("Failed to get gitlab project members: %v\n", err)
				continue
			}

			// Grant permissions to the user in sonar project
			for _, gitlabUser := range gitlabMembers {
				// Skip users not being synced this time
				if gitlabUser.Username != username {
					continue
				}
				fmt.Println("gitlabUser.Username:", gitlabUser.Username)
				sonarUser, e := mapSonarUsers[gitlabUser.Username]
				if e {
					// Grant permissions to users existing in both gitlab and sonar
					addData, err := AddUserPermissionToProject(sProject.Kee, sonarUser.Login, int(gitlabUser.AccessLevel))
					if err != nil {
						log.Logger.Warningf("Failed to grant gitlab project member permission: %v\n", err)
						continue
					}
					data = append(data, addData...)
				}

			}
		}

	}
	return
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func GetAllGroupName(groupName string) (subGroupNames []string) {
	splitGroupNames := strings.Split(groupName, ".")
	var subGroupName string
	for i := 0; i < len(splitGroupNames); i++ {
		if i == 0 {
			subGroupName = splitGroupNames[i]
		} else {
			subGroupName = subGroupName + "." + splitGroupNames[i]
		}
		subGroupNames = append(subGroupNames, subGroupName)
	}
	return subGroupNames
}

func AddUserPermissionToProject(projectKey, login string, level int) (data []PermissionData, err error) {
	var permissions []string
	var removePermissions []string
	if level == 30 {
		permissions = []string{"codeviewer", "user"}
		removePermissions = []string{"admin", "scan", "issueadmin", "securityhotspotadmin"}
	} else if level >= 40 {
		permissions = []string{"admin", "codeviewer", "issueadmin", "securityhotspotadmin", "user"}
		removePermissions = []string{"scan"}
		//removePermissions = []string{}
	} else {
		permissions = []string{}
		removePermissions = []string{"admin", "codeviewer", "issueadmin", "securityhotspotadmin", "user", "scan"}
	}

	for _, p := range permissions {
		fmt.Printf("add project: [%s], user: [%s], permission: [%s] \n", projectKey, login, p)
		data = append(data, PermissionData{Method: "add", User: login, Permission: p, Project: projectKey})
		err = sonarqube.AddUsersToProject(projectKey, login, p)
		if err != nil {
			log.Logger.Warningf("Failed to add sonar project member permission: %v\n", err)
			continue
		}
		fmt.Println("success!")
		time.Sleep(5 * time.Millisecond)
	}

	for _, p := range removePermissions {
		fmt.Printf("remove project: [%s], user: [%s], permission: [%s] \n", projectKey, login, p)
		data = append(data, PermissionData{Method: "remove", User: login, Permission: p, Project: projectKey})
		err = sonarqube.RemoveUsersToProject(projectKey, login, p)
		if err != nil {
			log.Logger.Warningf("Failed to remove sonar project member permission: %v\n", err)
			continue
		}
		fmt.Println("success!")
		time.Sleep(5 * time.Millisecond)
	}
	return
}

func UpdateGroupPermissionToProject(projectKey string) (data []PermissionData, err error) {
	// group permission update
	var allPermission = []string{"admin", "codeviewer", "issueadmin", "securityhotspotadmin", "user", "scan"}
	var adminPermission = []string{"admin", "user"}
	var adminRemovePermission = []string{"codeviewer", "issueadmin", "securityhotspotadmin", "scan"}
	data = []PermissionData{}

	permissionData, err := UpdateGroupPermission(projectKey, "sonar-users", []string{}, allPermission)
	if err != nil {
		log.Logger.Warningf("Failed to update sonar project sonar-users permissions: %v\n", err)
	}
	data = append(data, permissionData...)

	permissionData, err = UpdateGroupPermission(projectKey, "sonar-administrators", adminPermission, adminRemovePermission)
	if err != nil {
		log.Logger.Warningf("Failed to update sonar project sonar-administrators permissions: %v\n", err)
	}
	data = append(data, permissionData...)
	return
}

func RemoveUserCertainPermission(projectKey, login string, removePermissions []string) (data []PermissionData, err error) {
	for _, p := range removePermissions {
		fmt.Printf("remove project: [%s], user: [%s], permission: [%s] \n", projectKey, login, p)
		data = append(data, PermissionData{Method: "remove", User: login, Permission: p, Project: projectKey})
		err = sonarqube.RemoveUsersToProject(projectKey, login, p)
		if err != nil {
			log.Logger.Warningf("Failed to remove sonar project member permission: %v\n", err)
			continue
		}
	}
	return
}

func UpdateGroupPermission(projectKey, groupName string, addPermissions, removePermissions []string) (data []PermissionData, err error) {
	fmt.Println("sonar.groupName:", groupName)
	for _, permission := range removePermissions {
		fmt.Printf("remove project: [%s], groupName: [%s], permission: [%s] \n", projectKey, groupName, permission)
		data = append(data, PermissionData{Method: "remove", User: groupName, Permission: permission, Project: projectKey})
		err = sonarqube.RemoveGroupToProject(projectKey, groupName, permission)
		if err != nil {
			log.Logger.Warningf("Failed to remove sonar project group permission: %v\n", err)
			continue
		}
		fmt.Println("success!")
		time.Sleep(5 * time.Millisecond)
	}
	for _, permission := range addPermissions {
		fmt.Printf("add project: [%s], groupName: [%s], permission: [%s] \n", projectKey, groupName, permission)
		data = append(data, PermissionData{Method: "add", User: groupName, Permission: permission, Project: projectKey})
		err = sonarqube.AddGroupToProject(projectKey, groupName, permission)
		if err != nil {
			log.Logger.Warningf("Failed to add sonar project group permission: %v\n", err)
			continue
		}
		fmt.Println("success!")
		time.Sleep(5 * time.Millisecond)
	}
	return
}

type PermissionData struct {
	Method     string `json:"Method"`
	User       string `json:"User"`
	Permission string `json:"Permission"`
	Project    string `json:"Project"`
}
