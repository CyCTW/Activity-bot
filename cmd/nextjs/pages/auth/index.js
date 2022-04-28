import { Box, Button, VStack, Link } from "@chakra-ui/react";
import React, { useState, useEffect } from "react";

const Auth = (props) => {
  const { liff, liffError } = props;
  const [loginUser, setLoginUser] = useState();
  const [notifyLink, setNotifyLink] = useState();
// const [p1, setP1] = useState();
  const JsonToQuerystring = (payload) => {
    return (
      "?" +
      Object.keys(payload)
        .map((key) => {
          return (
            encodeURIComponent(key) + "=" + encodeURIComponent(payload[key])
          );
        })
        .join("&")
    );
  };
  const createNotifyLink = () => {
    let payload = {
      response_type: "code",
      client_id: process.env.LINE_NOTIFY_CLIENT_ID,
      redirect_uri: process.env.LINE_NOTIFY_REDIRECT_URI,
      scope: "notify",
      state: `${loginUser?.userId}_${loginUser?.displayName}`,
    };
    // setP1(payload)
    let notifyLink_ =
      "https://notify-bot.line.me/oauth/authorize" + JsonToQuerystring(payload);
    return notifyLink_;
  };
  useEffect(() => {
    // Get User Profile
    const getUserProfile = async () => {
      try {
        const profile = await liff.getProfile();
        setLoginUser(profile);
        console.log(profile)
      } catch (err) {
        console.log(err);
      }

    };
    if (liff !== null) {
      getUserProfile();
    }
  }, [liff]);

  useEffect(() => {
    if (loginUser !== null) {
        const link = createNotifyLink()
        setNotifyLink(link)
    }
  }, [loginUser])
  return (
    <div>
      <VStack>
        <h1>請點選以下開啟通知按鈕</h1>
        {loginUser && notifyLink && (
          <Box as='button' borderRadius='md' bg='tomato' color='white' px={4} h={8}>
            <Link href={notifyLink}>開啟通知</Link>
          </Box>
        )}
      </VStack>
    </div>
  );
};

export default Auth;
