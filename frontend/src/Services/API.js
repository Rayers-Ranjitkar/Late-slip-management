import axios from "axios";

const apiInstance = axios.create({
    baseURL:"http://localhost:8000",
});

//Post register data to API
export const postSignUpDataAPI = (userName,email,password) =>{
    return apiInstance.post("/admin/register",{
        fullname:userName,
        email:email,
        password:password
    });
}

//Post Login data to API
export const postLoginDataAPI = (email,password) =>{
    return apiInstance.post("/admin/login",{
        email:email,
        password:password
    });
}