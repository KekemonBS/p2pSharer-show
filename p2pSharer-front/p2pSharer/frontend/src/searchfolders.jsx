import {useState} from "preact/hooks";                           
import classes from "./searchfolders.css"
const SearchFolders = ({ no, Name, Size, FileTree, MagnetLink }) => {
    return (
        <div style="overflow-y: auto">
           <table>
               <thead>
                   <tr>
                       <th style="width: 3%">{no}</th>
                       <th title={Name}>{Name}</th>
                       <th style="width: 3%">{Size}</th>
                   </tr>
               </thead>
               <tbody>
                   <tr>
                       <td colspan="3" id="#tree">
                           <div id="filetree" style={{textAlign: "left"}}>
                               {FileTree}
                           </div>
                       </td>
                   </tr>
                   <tr>
                       <td style="width: 3%;text-align:center;"> 
                           <i class="fa fa-solid fa-magnet"></i>
                       </td>
                       <td colspan="2">
                           <a href={MagnetLink} target="_blank">{MagnetLink}</a>
                       </td>
                   </tr>
               </tbody>
           </table>
        </div>
  );
};

export default SearchFolders;
                                         
                                                                 

