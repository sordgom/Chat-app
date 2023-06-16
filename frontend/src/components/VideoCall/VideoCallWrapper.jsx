import { useLocation } from 'react-router-dom';
import VideoCall from './VideoCall';

function VideoCallWrapper() {
  const location = useLocation();
  
  // Pass location to the class-based component
  return <VideoCall location={location} />;
}

export default VideoCallWrapper;