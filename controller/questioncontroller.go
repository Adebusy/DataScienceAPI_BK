package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	cr "github.com/Adebusy/dataScienceAPI/crudal"
	"github.com/Adebusy/dataScienceAPI/model"
	ut "github.com/Adebusy/dataScienceAPI/utilities"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/swaggo/swag/example/celler/httputil"
)

// CreateNewQuestion godoc
// @Summary create new question
// @Produce json
// @Param user body model.Question true "create new question"
// @Success 200 {object} utilities.ResponseManager
// @Router /question/CreateNewQuestion/ [post]
func CreateNewQuestion(ctx *gin.Context) {
	var questionOBJ model.Question
	var resp ut.ResponseManager
	resp.ResponseCode = ""
	resp.ResponseDescription = ""
	if err := ctx.ShouldBindJSON(&questionOBJ); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	//validate Request
	validateRequestResp := ut.ValidateQuestionReq(questionOBJ)
	if validateRequestResp.ResponseCode != "" {
		ctx.JSON(http.StatusBadRequest, validateRequestResp)
		return
	}
	//check course name
	if CourseNameValidation := cr.CheckIfCourseExistBool(strings.ToUpper(questionOBJ.CourseName)); CourseNameValidation == false {
		resp.ResponseCode = "01"
		resp.ResponseDescription = "Course name supplied does not exist."
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// create question for the course
	if cr.CreateQuestion(questionOBJ) {
		resp.ResponseCode = "00"
		resp.ResponseDescription = fmt.Sprintf("Question created successfully for %s", questionOBJ.CourseName)
		ctx.JSON(http.StatusOK, resp)
	}
}

// FetchQuestionsByCourse godoc
// @Summary Fetch question for quis
// @Produce json
// @Param user body model.Question true "Fetch question for quis"
// @Success 200 {object} model.QuisRequest
// @Router /question/FetchQuestionsByCourse/{StudentID}/{CourseName} [get]
func FetchQuestionsByCourse(ctx *gin.Context) {
	var resp ut.ResponseManager
	resp.ResponseCode = ""
	resp.ResponseDescription = ""
	if courserName := ctx.Param("CourseName"); courserName == "" {
		resp.ResponseCode = "01"
		resp.ResponseDescription = "CourseName is required"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if courserName := ctx.Param("StudentID"); courserName == "" {
		resp.ResponseCode = "01"
		resp.ResponseDescription = "StudentID is required"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	//validate email address supplied
	if checkEmail := ut.ValidateEmail(ctx.Param("StudentID")); checkEmail == false {
		resp.ResponseCode = "01"
		resp.ResponseDescription = "Student Email address must be valid."
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	//use email fetch student details
	doStudentCheck, err := cr.GetStudentByEmailAddress(strings.ToUpper(ctx.Param("StudentID")))

	if err != nil {
		fmt.Println(err.Error)
		log.Panic(err)
	}
	//use code code fetch course detail
	doCheck := cr.CheckIfCourseExist(strings.ToUpper(ctx.Param("CourseName")))
	//confirm course name exist
	if doCheck.CourseCode == "" {
		resp.ResponseCode = "02"
		resp.ResponseDescription = "Course does not exist."
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	//fetchQuestion use course name and maxinum exam question per quis for this course
	CountValue, _ := strconv.Atoi(doCheck.QuestionCount)
	myquestion := cr.GetQuestion(doCheck.CourseName, CountValue)
	var res model.QuisRequest
	res.MyQuestion = myquestion
	res.CourseDetails = doCheck
	res.StudentDetails = doStudentCheck
	ctx.JSON(http.StatusOK, res)
}
