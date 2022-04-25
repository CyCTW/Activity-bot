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
	// stime := time.Date(2022, time.April, 22, 7, 19, 0, 0, time.UTC)
	activityTime := activity.Date
	remindTime := activityTime.Add(-time.Second * 30)

	activity_id := strconv.FormatUint(uint64(activity.ID), 10)
	message := fmt.Sprintf("您的活動 %v 將在一小時後舉辦", activity.Name)
	job, _ := s.Every(1).Day().StartAt(remindTime).Tag(activity_id).Tag(user.LineUserID).Do(NotifyUser, user.AccessToken, message)
	job.LimitRunsTo(1)
	// s.Every(5).Seconds().Do(task)
	s.StartAsync()
	log.Print("Scheduler end")
}

func RemoveScheduler(user *models.User, activityID string) {
	userLineID := user.LineUserID
	s.RemoveByTags(userLineID, activityID)
}

// aveIhis01dqAHEXA2AlDUZqY4R2m7nAuymDsyIl4rp3
