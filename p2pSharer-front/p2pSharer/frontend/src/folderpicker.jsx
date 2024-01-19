import { useState } from 'preact/hooks';
import {ShowFolderPicker} from "../wailsjs/go/main/App.js";

function FolderPicker() {
    const [selectedFolder, setSelectedFolder] = useState('');
    const divStyle = {
        display: "inline-block",
        height: "auto",
        width: "auto",
    };
    const pStyle = {
        display: "lnline-block",
        height: "100%",
        width: "auto",
        color: "white",
    };


    const selectFolder = () => {
        ShowFolderPicker("TEST").then((result) => (
            setSelectedFolder(result)
        ));
    };

    return (
        <div style={divStyle}>
            <button onClick={selectFolder}>
                Choose Folder</button>
            {
                selectedFolder && 
                <p style={pStyle}>
                    Selected Folder: {selectedFolder}
                </p>
            }
        </div>
    );
}

export default FolderPicker;
