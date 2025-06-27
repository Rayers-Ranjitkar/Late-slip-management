import Authentication from "../Components/Authentication/Authentication";
import {postSignUpDataAPI} from "../Services/API";
import { toast } from "react-toastify";
import { useNavigate } from "react-router-dom";

const Registration = () => {
  const navigate = useNavigate();
  const onRegistrationFormSubmit = async (userName,email, password) => {
    try {
      const resPostSignUpDataAPI = await postSignUpDataAPI(userName,email, password);
      console.log("resPostSignUpDataAPI", resPostSignUpDataAPI);
      toast.success("User Registered in successfully");
      navigate("/login");
    } catch (error) {
      console.log("Failed SignUp. Error -> ", error);
      toast.error("Failed SignUp. User already exists!", error);
    }
  };

  return (
    <>
      <Authentication
        isLoginPage={false}
        onFormSubmit={onRegistrationFormSubmit}
      />
    </>
  );
};

export default Registration;
