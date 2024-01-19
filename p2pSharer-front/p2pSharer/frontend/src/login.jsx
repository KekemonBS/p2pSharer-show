import {useState} from "preact/hooks";
import {POSTGo} from "../wailsjs/go/main/App.js";

function Login(props) {
 let location = props.location
 const [username, setUsername] = useState('');
 const [password, setPassword] = useState('');
 const [isLoggedIn, setIsLoggedIn] = useState(false);

 const inputStyles = {
    display: "block",
    marginBottom: "1vh",
    width: "100%",
 };

 const join = async () => {
    let reqString = `${location}/api/v1/tracker/init?username=${username}&password=${password}`;
    console.log(reqString)
    POSTGo(reqString).then((response) => {
        if (response == " 200 OK") {
            alert("HTTP Err : " + response);
        }
    });
 }
    
 const leave = async () => {
    let reqString = `${location}/api/v1/tracker/leave?username=${username}&password=${password}`;
    console.log(reqString)
    POSTGo(reqString).then((response) => {
        console.log(response)
        if (response == " 200 OK") {
            alert("HTTP Err : " + response);
        }
    });
 }

 const handleJoin = () => {
    if (!isLoggedIn) {
        join();
    } else {
        leave();
    }
    setIsLoggedIn(!isLoggedIn);
 };



 return (
   <div>
     <input
       style={inputStyles}
       type="text"
       placeholder="Username"
       value={username}
       onInput={(e) => setUsername(e.target.value)}
     />
     <input
       style={inputStyles}
       type="password"
       placeholder="Password"
       value={password}
       onInput={(e) => setPassword(e.target.value)}
     />
     <button onClick={handleJoin}>
       {isLoggedIn ? 'Leave' : 'Join'}
     </button>
   </div>
 );
}

export default Login;

