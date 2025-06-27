/* API Import */
import { postLoginDataAPI } from "../Services/API";

/* Other Imports */
import {toast} from 'react-toastify';
import Authentication from "../Components/Authentication/Authentication";
import { useNavigate } from "react-router-dom";

const Login = () => {
    const navigate = useNavigate(); // when we call useNavigate() function it returns a function (let’s call it a “navigate function”). // 'navigate' variable is now a function and inside that when we pass a path navigate('path') it redirects us to that path hehe
     //console.log(localStorage.getItem('token'));
  /* Login API Request */
  const onLoginSubmit = async (email,password) => {
    try {
      const resPostLoginDataAPI = await postLoginDataAPI(email, password);
      console.log("resPostLoginDataAPI", resPostLoginDataAPI);
      toast.success("User logged in successfully");

      /* After the user is logged in */
      localStorage.setItem('token',resPostLoginDataAPI.data.token);
      console.log("User logged in successfully + Stored token in localStorage");
      navigate('/adminDashboard');
      
    } catch (error) {
      console.log("Failed Login ->", error);
      toast.error("Invalid Email or password!", error);
    }
  };

  return (
    <>
      <Authentication isLoginPage={true} onFormSubmit={onLoginSubmit} />
    </>
  );
};

export default Login;
