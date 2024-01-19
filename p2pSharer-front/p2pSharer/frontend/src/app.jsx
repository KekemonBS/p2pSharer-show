import './app.css';
import {useState} from "preact/hooks";
import {Fragment} from "preact";

import Login from './login.jsx';
import FolderPicker from './folderpicker.jsx';
import SearchFolders from './searchfolders.jsx';

//import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
//import { faCoffee } from "@fortawesome/free-solid-svg-icons";

export function App(props) {
    const location = "http://localhost:8080"
    let resultsPresent = true;
    let searchResult = [
        {
            Name: "Rimworld-jc141",
            Size: "280.08 MiB",
            FileTree: `Rimworld-jc141/files/groot.dwarfsi
Rimworld-jc141/files/logo.txt.gz
Rimworld-jc141/settings.sh
Rimworld-jc141/start.n.sh
Rimworld-jc141/files
Rimworld-jc141`,
            MagnetLink: "Qma4XLPkrRNwKaqEXUTx3HyUq1ymZ8f8Kns2jvczSsuGvn",
        },  
    ];


    const [hosterUsername, setHosterUsername] = useState('');
    const [hosterPassword, setHosterPassword] = useState('');
    
    const divStyles = {
        display: "inline-block",
        width: "46vw",
        padding: "2vw",
    };
   
    const inputStyles = {
        display: "inline-block",
        marginBottom: "1vh",
        width: "100%",
    };

    const [isResultsVisible, setIsResultsVisible] = useState(false);
    const resultsStyles = {
        visibility: !isResultsVisible ? "visible" : "hidden",
        display: "inline-block",
        overflow: "auto",
        height: "68vh",
        border: "2px solid black",
        backgroundColor: "grey",

    };   

    const showResults = () => {
        if (!isResultsVisible) {
            setIsResultsVisible(true);
        } else {
            setIsResultsVisible(false);
        }
    };
    
    const searchButtonStyles = {
        marginBottom: "1vh",
    };

    const hostFolder = () => {
        ShowFolderPicker(hosterUsername, hosterPassword).then((result) => (
            setSelectedFolder(result)
        ));
    };

    const searchHosted = async () => {
        if (query.trim() == "") 
            return;
        let reqString = `${location}/api/v1/tracker/read?${hosterUsername}&${hosterPassword}`
        let response = await fetch(reqString);
        console.log(`reqString : ${reqString}`);
        console.log(`code : ${response.status}`);
        if (response.ok) {
            let respObj = await response.json();
            resultsFound = parseInt(respObj.Total, 10);
            searchResult =  respObj.Results;
            if (searchResult != "") {
                resultsPresent = true;
            } else {
                alert("Nothing found.");
            }
        } else {
            alert("HTTP Err : " + response.status);
        }
    }



    return (
        <>
            <div id="App">
                <div style={divStyles}> 
                    <Login location={location}/>
                    <FolderPicker/>
                </div>

                <div style={divStyles}>
                      <div style="owerflow: hidden;">
                      <input
                        style={inputStyles}
                        type="text"
                        placeholder="Hoster Username"
                        value={hosterUsername}
                        onInput={(e) => setHosterUsername(e.target.value)}
                      />
                      <input
                        style={inputStyles}
                        type="password"
                        placeholder="Hoster Password"
                        value={hosterPassword}
                        onInput={(e) => setHosterPassword(e.target.value)}
                      />
                        
                      <button style={searchButtonStyles} onClick={showResults}>
                       Search 
                      </button>
                      </div>

                      <div style={resultsStyles}>
                          {resultsPresent && (
                              searchResult.map((d, i) => (
                                  <Fragment key={i + 1}>
                                  <div style="padding: 10px">
                                      <SearchFolders
                                       no={i + 1}
                                       Name={d.Name}
                                       Size={d.Size}
                                       FileTree={d.FileTree}
                                       MagnetLink={d.MagnetLink}
                                      />
                                  </div>
                                  </Fragment>
                              ))
                          )}
                      </div>

                </div>
            </div>
        </>
    )
}


