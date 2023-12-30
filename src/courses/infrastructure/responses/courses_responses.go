package responses

import (
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"
)

type EnrolledCourse struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func GetEnrolledCourseFromCourseEntity(course *entities.Course) EnrolledCourse {
	return EnrolledCourse{
		UUID:  course.UUID,
		Name:  course.Name,
		Color: course.Color,
	}
}

type EnrolledCoursesResponse struct {
	Courses       []EnrolledCourse `json:"courses"`
	HiddenCourses []EnrolledCourse `json:"hidden_courses"`
}

func GetResponseFromDTO(dto *dtos.EnrolledCoursesDto) *EnrolledCoursesResponse {
	enrolledCourses := []EnrolledCourse{}
	hiddenCourses := []EnrolledCourse{}

	for _, course := range dto.Courses {
		enrolledCourses = append(enrolledCourses, GetEnrolledCourseFromCourseEntity(&course))
	}

	for _, course := range dto.HiddenCourses {
		hiddenCourses = append(hiddenCourses, GetEnrolledCourseFromCourseEntity(&course))
	}

	return &EnrolledCoursesResponse{
		Courses:       enrolledCourses,
		HiddenCourses: hiddenCourses,
	}
}
