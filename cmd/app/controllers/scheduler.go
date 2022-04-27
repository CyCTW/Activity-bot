package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cyctw/line-profile-bot/cmd/app/models"
	"github.com/go-co-op/gocron"
)

var s = gocron.NewScheduler(time.UTC)

func StartScheduler(user *models.User, activity *models.Activity) {
	activityTime := activity.Date

	// Note: Default set to alert one minute before
	remindTime := activityTime.Add(-time.Minute)

	activity_id := strconv.FormatUint(uint64(activity.ID), 10)
	message := fmt.Sprintf("您的活動 %v 將在一分鐘後舉辦", activity.Name)
	job, _ := s.Every(1).Day().StartAt(remindTime).Tag(activity_id).Tag(user.LineUserID).Do(NotifyUser, user.AccessToken, message)
	job.LimitRunsTo(1)

	s.StartAsync()
	log.Print("Scheduler end")
}

func RemoveScheduler(user *models.User, activityID string) {
	userLineID := user.LineUserID
	s.RemoveByTags(userLineID, activityID)
}
