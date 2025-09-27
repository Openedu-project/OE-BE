package guards

type Role string

const (
	RoleSysAdmin Role = "sysadmin"
	RoleAdmin    Role = "admin"
	RoleCreator  Role = "creator"
	RoleLearner  Role = "learner"
)

const (
	BlogCreate  Permission = "blog.create"
	BlogRead    Permission = "blog.read"
	BlogUpdate  Permission = "blog.update"
	BlogDelete  Permission = "blog.delete"
	BlogPublish Permission = "blog.publish"
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

	PermEnrollInCourse Permission = "course.enroll"

	PermViewMyCertificates Permission = "certificate.view.my"
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
		BlogCreate,
		BlogRead,
		BlogUpdate,
		BlogDelete,
		BlogPublish,
	},
	RoleAdmin: {
		PermUserManage,
		PermCourseCRUD,
		PermCourseView,
		PermContentCRUD,
		PermContentView,
		PermReportView,
		BlogCreate,
		BlogRead,
		BlogUpdate,
		BlogDelete,
		BlogPublish,
	},
	RoleCreator: {
		PermCourseCRUD,
		PermCourseView,
		PermContentCRUD,
		PermContentView,
		PermReportView,
		BlogCreate,
		BlogRead,
		BlogPublish,
	},
	RoleLearner: {
		PermCourseView,
		PermContentView,
		PermReportView,
		BlogRead,
		PermEnrollInCourse,
		PermViewMyCertificates,
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
