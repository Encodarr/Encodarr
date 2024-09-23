import styles from "./Media.module.scss";
import Footer from "../footer/Footer";
import { useContext, useState } from "react";
import { WebSocketContext } from "../../contexts/webSocketContext";
import MediaModel from "../modals/mediaModal/MediaModal";
import Posters from "../posters/Posters";
import Overviews from "../overviews/Overviews";
import MediaToolbar from "../toolbars/mediaToolbar/MediaToolbar";
import sortAndFilter from "../../utils/sortAndFilter";
import MediaTable from "../tables/mediaTable/MediaTable";

const Media = () => {
  const wsContext = useContext(WebSocketContext);
  const movies = wsContext?.data?.movies;
  const series = wsContext?.data?.series;
  const settings: any = wsContext?.data?.settings
  ? Object.keys(wsContext?.data?.settings).reduce((acc, key) => {
      acc[key] = wsContext?.data?.settings[key].value;
      return acc;
    }, {})
  : {};
  const profiles = wsContext?.data?.profiles;
  const view = settings?.mediaView
  const sort = settings?.mediaSort
  const filter = settings?.mediaFilter;
  const sortDirection = settings?.mediaSortDirection;
  const sortedMedia = sortAndFilter(
    series,
    movies,
    profiles,
    sort,
    sortDirection,
    filter
  );
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [content, setContent] = useState({});
  const [selected, setSelected] = useState(null);
  return (
    <div className={styles.media}>
      <MediaToolbar
        selected={selected}
        setSelected={setSelected}
        setContent={setContent}
        setIsModalOpen={setIsModalOpen}
        settings={settings}
        view={view}
        system={wsContext?.data.system}
      />
      <MediaModel
        type={view}
        isOpen={isModalOpen}
        setIsOpen={setIsModalOpen}
        content={content}
        setContent={setContent}
      />
      <div className={styles.mediaContent}>
        <div className={styles.contentContainer}>
          {view === "table" && (
            <MediaTable
              settings={settings}
              profiles={profiles}
              sortedMedia={sortedMedia}
            />
          )}
          {view === "posters" && (
            <Posters settings={settings} sortedMedia={sortedMedia || []} profiles={profiles} />
          )}
          {view === "overview" && (
            <Overviews
              sortedMedia={sortedMedia}
              settings={settings}
              profiles={profiles}
            />
          )}
          <div className={styles.footerContent}>{series && <Footer />}</div>
        </div>
      </div>
    </div>
  );
};
export default Media;
