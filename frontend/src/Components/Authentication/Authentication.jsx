/* Hooks Import */
import { useState } from "react";
import { useNavigate } from "react-router-dom";

/* Icons & Img Import */
import IconEmail from "../../assets/AuthenticationIcons/IconEmail.png";
import IconPassword from "../../assets/AuthenticationIcons/IconPassword.png";
import { FaUser } from "react-icons/fa";

const Authentication = ({ isLoginPage, onFormSubmit }) => {
  const navigate = useNavigate();
  //const [isLoginPage, setIsLoginPage] = useState(true); //Just for test

  const [email, setNewEmail] = useState();
  const [password, setNewpassword] = useState();
  const [username, setUsername] = useState(); //For signUp page

  const [isTypingPassword, setIsTypingPassword] = useState(false); // default: show

  /* Handling input changes */
  const handlePasswordChange = (e) => {
    setNewpassword(e.target.value);
  };

  const handleEmailChange = (e) => {
    setNewEmail(e.target.value);
  };

  const handleUsernameInput = (e) => {
    setUsername(e.target.value);
  };

  /* Showing the visibility of the password icon */
  const handlePasswordFocus = () => {
    setIsTypingPassword(true);
  };

  const handlePasswordBlur = () => {
    setIsTypingPassword(false);
  };

  return (
    <>
      <form className="font-quicksand">
        {/* Outermost div */}
        <div className="min-h-screen flex items-center  flex-col gap-4 justify-center px-20">
          {/* Upper Div */}
          <div className=" relative flex border-1 border-[#74c044] flex-wrap flex-col gap-4 min-w-[80vw] sm:min-w-[50%] 2xl:min-w-[40%] px-6 p-2 py-6 items-center ">
            <h1 className="font-quicksand font-bold text-center text-2xl mb-6 min-w-[90%] 2xl:text-4xl">
              Late Slip Management
            </h1>

            {/* Full name Input Div*/}
            {!isLoginPage && (
              <div className=" min-w-[90%] sm:min-w-[45%] 2xl:min-w-[60%] relative ">
                <input
                  placeholder="Username"
                  className="bg-[#DCE0E3] text-black font-bold text-xs 2xl:text-base h-10 2xl:h-12 w-full rounded-lg p-4"
                  type=""
                  onChange={handleUsernameInput}
                />
                {/* Password Img span */}
                <span className="absolute z-10 top-[20%] right-3.5 2xl:right-3.5 text-3xl text-[#6b6767] ">
                  <FaUser className="w-5 h-5" />
                </span>
              </div>
            )}

            {/* Email Input Div*/}
            <div className=" min-w-[90%] sm:min-w-[45%] 2xl:min-w-[60%] relative ">
              <input
                placeholder="Email"
                className="bg-[#DCE0E3] text-black font-bold text-xs 2xl:text-base h-10 2xl:h-12 w-full rounded-lg p-4"
                onChange={handleEmailChange}
              />
              {/* Email img span */}
              <span className="absolute z-10 top-[15%] right-3.5 2xl:right-3.5 text-3xl text-[#6b6767] ">
                <img src={IconEmail} alt="Email Icon" className="w-6 h-6" />
              </span>
            </div>

            {/* Password Input Div*/}
            <div className=" min-w-[90%] sm:min-w-[45%] 2xl:min-w-[60%] relative ">
              <input
                placeholder="Password"
                className="bg-[#DCE0E3] text-black font-bold text-xs 2xl:text-base h-10 2xl:h-12 w-full rounded-lg p-4"
                type="password"
                onChange={handlePasswordChange}
                onFocus={handlePasswordFocus}
                onBlur={handlePasswordBlur}
              />
              {/* Password Img span */}
              <span className="absolute z-10 top-[15%] right-3.5 2xl:right-3.5 text-3xl text-[#6b6767] ">
                {!isTypingPassword && (
                  <img
                    src={IconPassword}
                    alt="Password Icon"
                    className="w-6 h-6"
                  />
                )}
              </span>

              {/* Forgot pass Paragraph */}
              {isLoginPage && (
                <p className="mt-4.5 text-xs text-black font-semibold sm:text-sm 2xl:text-base cursor-pointer">
                  Forgot Password!
                </p>
              )}
            </div>

            {/* Submit Btn Div*/}
            <div
              onClick={() => {
                // if (!username || !email || !password) {
                //   toast.error("Please fill in all fields");
                //   return;
                // }
                if (isLoginPage) {
                  onFormSubmit(email, password);
                } else {
                  onFormSubmit(username, email, password); //Let's not overcomplicate for the project//we wanna clear up the fields so we can't call it here as we can't store the repo, we wanna store this reponse and if respose = success then, only we wanna clear up the fields so, we will have to create another form submit function and call it there
                }
              }}
              className="flex h-10 2xl:h-12 mt-3.5 bg-[#74c044] hover:bg-[#5e9e36] rounded-md justify-center cursor-pointer font-bold text-xs 2xl:text-base items-center min-w-[90%] sm:min-w-[45%] text-white 2xl:min-w-[60%]"
            >
              {isLoginPage ? "Login" : "Sign Up"}
            </div>
          </div>

          {/* Lower Div   */}
          <div className="min-h-[10vh]  min-w-[80vw]  sm:min-w-[50%] 2xl:min-w-[40%] px-6 flex justify-center items-center gap-0.5 sm:gap-4 2xl:gap-6 text-sm 2xl:text-lg border-1 border-[#74c044] ">
            <p className="font-bold font">
              {isLoginPage ? "Don't have an account?" : "Have an Account?"}
            </p>
            <button
              type="button"
              onClick={() => {
                isLoginPage ? navigate("/signUp") : navigate("/login");
              }} /* UseNavigate() is not a component and is a callback function. So, it is called after certain event to re-direct something. which <Navigate> component can't. Since, this is a callback func we will pass this inside arrow function so that it won't immediately be executed */
              className="flex w-20 h-8 2xl:h-12 2xl:w-35 bg-[#74c044] hover:bg-[#5e9e36] text-white  rounded-md 2xl:rounded-2xl justify-center cursor-pointer font-bold text-xs 2xl:text-base items-center ml-5.5 sm:ml-0"
            >
              {isLoginPage ? "Sign Up" : "Login"}
            </button>
          </div>
        </div>
      </form>
    </>
  );
};

export default Authentication;
