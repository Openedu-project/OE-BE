package guards

type Role string

const (
	RoleSysAdmin Role = "sysadmin"
	RoleAdmin    Role = "admin"
	RoleCreator  Role = "creator"
	RoleLearner  Role = "learner"
)

type Permission string

const (
	// User
	PermUserManage Permission = "user.manage"

	// Course
	PermCourseCRUD Permission = "course.crud"
	PermCourseView Permission = "course.view"

	// Content
	PermContentCRUD Permission = "content.crud"
	PermContentView Permission = "content.view"

	// Reports
	PermReportView Permission = "report.view"

	// System
	PermSystemConfig Permission = "system.config"
)

var rolePermissions = map[Role][]Permission{
	RoleSysAdmin: {
		PermUserManage,
		PermCourseCRUD,
		PermCourseView,
		PermContentCRUD,
		PermContentView,
		PermReportView,
		PermSystemConfig,
	},
	RoleAdmin: {
		PermUserManage,
		PermCourseCRUD,
		PermCourseView,
		PermContentCRUD,
		PermContentView,
		PermReportView,
	},
	RoleCreator: {
		PermCourseCRUD,
		PermCourseView,
		PermContentCRUD,
		PermContentView,
		PermReportView,
	},
	RoleLearner: {
		PermCourseView,
		PermContentView,
		PermReportView,
	},
}

func HasPermission(role Role, perm Permission) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}
