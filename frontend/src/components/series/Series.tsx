import { useContext, useEffect, useRef, useState } from "react";
import styles from "./Series.module.scss";
import Drive from "../svgs/hard_drive.svg?react";
import Profile from "../svgs/person.svg?react";
import Monitored from "../svgs/bookmark_filled.svg?react";
import Unmonitored from "../svgs/bookmark_unfilled.svg?react";
import Continuing from "../svgs/play_arrow.svg?react";
import Ended from "../svgs/stop.svg?react";
import Network from "../svgs/tower.svg?react";
import Season from "../season/Season";
import { WebSocketContext } from "../../contexts/webSocketContext";
import SeriesModal from "../modals/seriesModal/SeriesModal";
import SeriesToolbar from "../toolbars/seriesToolbar/SeriesToolbar";
import { formatSize } from "../../utils/format";
import FolderIcon from "../svgs/folder.svg?react";

const Series = ({ seriesName }: any) => {
  const wsContext = useContext(WebSocketContext);
  const profiles = wsContext?.data?.profiles;
  const series: any = 
  wsContext?.data?.series && profiles
    ? wsContext?.data?.series.find((s: any) => s.id === seriesName)
    : {};
  const system: any = wsContext?.data?.system
    ? Object.keys(wsContext?.data?.system).reduce((acc, key) => {
      acc[key] = wsContext?.data?.system[key].value;
      return acc;
    }, {})
    : {};
  const [content, setContent] = useState<any>({});
  const handleEditClick = () => {
    setIsModalOpen(true);
    setContent(series);
  };
  const [selected, setSelected] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const status = series?.status;
  const network = series?.networks;
  const genre = series?.genre;
  const firstAirDate = series?.releaseDate?.split("-")[0].trim();
  const lastAirDate = series?.lastAirDate?.split("-")[0].trim();
  const overview = series?.overview;
  const runYears =
    status === "Ended" ? firstAirDate + "-" + lastAirDate : firstAirDate + "-";
  const [backdropSrc, setBackdropSrc] = useState<string | null>("");
  const [posterSrc, setPosterSrc] = useState<string | null>("");
  const loaded = useRef(false);
  useEffect(() => {
    if (loaded.current == true) {
      return;
    }
    const fetchImage = async (
      path: string,
      setSrc: (src: string | null) => void
    ) => {
      try {
        let cache = null;
        let cachedResponse = null;
        if ("caches" in window) {
          cache = await caches.open("image-cache");
          cachedResponse = await cache.match(
            `/api/series/${series?.id}/${path}`
          );
        }

        if (cachedResponse) {
          const blob = await cachedResponse.blob();
          setSrc(URL.createObjectURL(blob));
        } else {
          const response = await fetch(`/api/series/${series?.id}/${path}`, {
            headers: {
              Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
          });
          if (response.status !== 200) {
            setSrc(null);
            return;
          }
          const clonedResponse = response.clone();
          const blob = await response.blob();
          setSrc(URL.createObjectURL(blob));
          if (cache) {
            cache.put(`/api/series/${series?.id}/${path}`, clonedResponse);
          }
        }
      } catch (e) {
        console.log(e);
      }
    };

    if (series?.id && series?.id !== "") {
      fetchImage("backdrop", setBackdropSrc);
      fetchImage("poster", setPosterSrc);
      loaded.current = true;
    }
  }, [series?.id]);

  const profileName = profiles?.find((profile: any) => profile.id === series?.profileId)?.name || "";
  return (
    <div className={styles.series}>
      <SeriesToolbar
        series={series}
        system={system}
        selected={selected}
        setSelected={setSelected}
        handleEditClick={handleEditClick}
        seriesName={seriesName}
      />
      <SeriesModal
        isOpen={isModalOpen}
        setIsOpen={setIsModalOpen}
        content={content}
        setContent={setContent}
        profiles={profiles}
      />
      <div className={styles.seriesContent}>
        <div className={styles.header}>
          <img
            className={styles.backdrop}
            src={backdropSrc || "/backdrop.jpg"}
            alt="backdrop"
          />
          <div className={styles.filter} />
          <div className={styles.content}>
            <img
              className={styles.poster}
              src={posterSrc || "/poster.png"}
              alt="poster"
            />
            <div className={styles.details}>
              <div className={styles.titleRow}>
                <div className={styles.headerIcon}>
                  {series?.monitored ? (
                    <Monitored className={styles.monitoredSVG} />
                  ) : (
                    <Unmonitored className={styles.monitoredSVG} />
                  )}
                </div>

                {series?.name ? series?.name : series?.id}
              </div>
              <div className={styles.seriesDetails}>
                <span className={styles.runtime}>
                  {series?.runtime ? series?.runtime : "-"} Minutes
                </span>
                {genre ? <span className={styles.genre}>{genre}</span> : <></>}
                {status ? (
                  <span className={styles.runYears}>{runYears}</span>
                ) : (
                  <></>
                )}
              </div>
              <div className={styles.tags}>
                <div className={styles.tag}>
                  <div className={styles.icon}>
                    <FolderIcon className={styles.svg} />
                  </div>
                  {"/series/" + series?.id}
                </div>

                <div className={styles.tag}>
                  <div className={styles.icon}>
                    <Drive className={styles.svg} />
                  </div>
                  {formatSize(series?.size)}
                </div>
                <div className={styles.tag}>
                  <div className={styles.icon}>
                    <Profile className={styles.svg} />
                  </div>
                  {profileName}
                </div>
                <div className={styles.tag}>
                  <div className={styles.icon}>
                    {series?.monitored ? (
                      <Monitored className={styles.svg} />
                    ) : (
                      <Unmonitored className={styles.svg} />
                    )}
                  </div>
                  {series?.monitored ? "Monitored" : "Unmonitored"}
                </div>
                {status ? (
                  <div className={styles.tag}>
                    <div className={styles.icon}>
                      {status === "Ended" ? (
                        <Ended className={styles.svg} />
                      ) : (
                        <Continuing className={styles.svg} />
                      )}
                    </div>
                    {status}
                  </div>
                ) : (
                  <></>
                )}
                {network ? (
                  <div className={styles.tag}>
                    <div className={styles.icon}>
                      <Network className={styles.svg} />
                    </div>
                    {network}
                  </div>
                ) : (
                  <></>
                )}
              </div>
              <div className={styles.overview}>{overview}</div>
            </div>
          </div>
        </div>
        <div className={styles.seasonsContainer}>
          {Object.values(series?.seasons || {})
            .reverse()
            .map((season: any) => {
              return (
                <Season
                  season={season}
                  monitored={series?.monitored}
                  key={season?.id}
                />
              );
            })}
        </div>
      </div>
    </div>
  );
};
export default Series;
