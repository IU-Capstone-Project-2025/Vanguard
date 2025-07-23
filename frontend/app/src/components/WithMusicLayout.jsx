import { Outlet } from "react-router-dom";
import BackgroundMusic from "./background/BackgroundMusic.jsx";

const WithMusicLayout = () => {
  return (
    <>
      <BackgroundMusic />
      <Outlet /> {/* Здесь будут PageOne и PageTwo */}
    </>
  );
};

export default WithMusicLayout;
