import { useRouter } from "next/router";
import React, { useState, useEffect } from "react";
// import { getActivity } from "../../../service";

const Activity = () => {
  const router = useRouter();
  const { activityID } = router.query;

  const [users, setUsers] = useState();
  const [activity, setActivity] = useState();

  // useEffect(() => {
  //   // Get activity data
  //   const getActivityData = async () => {
  //     const res = await getActivity(activityID);
  //     const users = res.data.users;
  //     const activity = res.data.activity;
  //     setUsers(users);
  //     setActivity(activity);
  //     console.log(users);
  //     console.log(activity);
  //   };

  //   getActivityData();
  // }, [activityID]);

  return (
    <div>
      {/* <p>{activity?.Name}</p>
      <p>{activity?.Date}</p>
      <p>{activity?.Place}</p>

      {users.map((idx, user) => {
        return <div key={idx}>{user?.Name}</div>;
      })} */}

      haha
      {activityID}
    </div>
  );
};

export default Activity;
