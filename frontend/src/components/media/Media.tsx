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
  const settings = wsContext?.data?.settings;
  const profiles = wsContext?.data?.profiles;
  const view = settings?.media_view.value;
  const sort = settings?.media_sort.value;
  const filter = settings?.media_filter.value;
  const sortDirection = settings?.media_sort_direction.value;
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
  console.log("sorted", wsContext?.data);
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
            <Posters settings={settings} sortedMedia={sortedMedia || []} />
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
