import axios from "axios";

const services = {};

const instance = axios.create({
  baseURL: "",
  withCredentials: true,
});

export const createActivity = ({ name, dateString, place, idToken }) => {
  return instance.post(`/activity`, {
    name,
    date: dateString,
    place,
    idToken
  });
};

export const getActivity = (activityID) => {
  return instance.get(`/activity/${activityID}`)
}