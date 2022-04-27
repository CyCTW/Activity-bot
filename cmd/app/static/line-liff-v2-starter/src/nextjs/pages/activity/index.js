import { Box, Button, VStack } from "@chakra-ui/react";
import { useRouter } from "next/router";
import React, { useState, useEffect } from "react";
import { getActivity } from "../../service";
// import { getActivity } from "../../../service";

const Activity = (props) => {
  const { liff, liffError } = props;

  const router = useRouter();
  console.log(router.asPath);

  const [loginUser, setLoginUser] = useState({});
  const [users, setUsers] = useState();
  const [activity, setActivity] = useState();

  useEffect(() => {
    const path = router.asPath;
    let pos = path.search("id=");
    const activityID = path.slice(pos + 3);
    console.log(activityID);

    // Get activity data
    const getActivityData = async () => {
      const res = await getActivity(activityID);
      const users = res.data.users;
      const activity = res.data.activity;
      setUsers(users);
      setActivity(activity);
      console.log(users);
      console.log(activity);
    };
    getActivityData();
  }, []);

  useEffect(() => {
    // Get User Profile
    const getUserProfile = async () => {
      try {
        const profile = await liff.getProfile();
        setLoginUser(profile);
        console.log("Profile")
        console.log(profile)
      } catch (err) {
        console.log(err);
      }
    };
    if (liff !== null) {
      getUserProfile();
    }
  }, [liff]);

  return (
    <div>
      <VStack>
        <h1 className="home__title">活動資訊</h1>
        <p>活動名稱: {activity?.Name}</p>
        <p>活動日期: {activity?.Date}</p>
        <p>活動地點: {activity?.Place}</p>
        <h3 className="home__title">參加者名單:</h3>

        {users &&
          users.map((user, idx) => {
            return <div key={idx}>{user?.Name}</div>;
          })}
        {users && users.some((item) => item.LineUserID == loginUser?.userId) ? (
          <Box>您已經參加此活動!</Box>
        ) : (
          <></>
        )}
      </VStack>
    </div>
  );
};

export default Activity;
