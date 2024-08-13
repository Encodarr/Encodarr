import styles from "./Poster.module.scss";
import { useContext, useEffect, useState } from "react";
import { WebSocketContext } from "../../contexts/webSocketContext";
import { Link } from "react-router-dom";
const PosterComponent = ({ media, posterWidth, posterHeight }: any) => {
  const wsContext = useContext(WebSocketContext);
  const profiles = wsContext?.data?.profiles;
  const settings = wsContext?.data?.settings;
  const type = media?.episode_count != undefined ? "series" : "movies";
  const progress = () => {
    if (media?.episode_count == undefined) {
      return media?.missing == true ? "0%" : "100%";
    }
    return media?.episode_count === 0
      ? "100%"
      : ((media?.episode_count - media?.missing_episodes) /
          media?.episode_count || 0) *
          100 +
          "%";
  };
  const backgroundColor = () => {
    if (progress() === "100%") {
      if (media?.missing_episodes == undefined) {
        return "rgb(39, 194, 76)";
      }
      return media?.status === "Ended"
        ? "rgb(39, 194, 76)"
        : "rgb(93, 156, 236)";
    } else {
      return media?.monitored ? "rgb(240, 80, 80)" : "rgb(255, 165, 0)";
    }
  };

  const [imgSrc, setImgSrc] = useState<string | null>("");
  useEffect(() => {
    const fetchImage = async () => {
      try {
        let cachedResponse = null;
        let cache = null;
        if ("caches" in window) {
          cache = await caches.open("image-cache");
          cachedResponse = await cache.match(
            `/api/poster/${type}/${media?.id}`
          );
        }

        if (cachedResponse) {
          const blob = await cachedResponse.blob();
          setImgSrc(URL.createObjectURL(blob));
        } else {
          const response = await fetch(`/api/poster/${type}/${media?.id}`, {
            headers: {
              Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
          });

          if (response.status !== 200) {
            setImgSrc(null);
            return;
          }

          const clonedResponse = response.clone();
          const blob = await response.blob();
          setImgSrc(URL.createObjectURL(blob));
          if (cache) {
            cache.put(`/api/poster/${type}/${media?.id}`, clonedResponse);
          }
        }
      } catch (e) {
        console.log(e);
      }
    };

    fetchImage();
  }, [media?.id, type]);
  return (
    <div
      className={styles.cardArea}
      style={{ maxWidth: posterWidth, maxHeight: posterHeight }}
    >
      <Link to={`${type}/${media?.id}`} className={styles.poster}>
        <div className={styles.card}>
          <div className={styles.cardContent}>
            <img
              className={styles.img}
              src={imgSrc || "/poster.png"}
              alt={media?.name}
              style={{ maxWidth: posterWidth, maxHeight: posterHeight }}
            />
            <div className={styles.footer}>
              <div className={styles.progressBar}>
                <div
                  className={styles.progress}
                  style={{
                    backgroundColor: backgroundColor(),
                    width: progress(),
                    height:
                      settings.media_poster_detailedProgressBar == "1"
                        ? "15px"
                        : "5px",
                  }}
                />
                {settings?.media_poster_detailedProgressBar == "1" && (
                  <div className={styles.detailText}>
                    {media?.episode_count == undefined ? (
                      <>{media?.missing ? "0/1" : "1/1"}</>
                    ) : (
                      <>
                        {media?.episode_count - media?.missing_episodes}/
                        {media?.episode_count}
                      </>
                    )}
                  </div>
                )}
              </div>
              {settings?.media_poster_showTitle == "1" && (
                <div className={styles.name}>
                  {media?.name ? media?.name : media?.id}
                </div>
              )}
              {settings?.media_poster_showMonitored == "1" && (
                <div className={styles.status}>
                  {media?.monitored ? "Monitored" : "Unmonitored"}
                </div>
              )}
              {settings?.media_poster_showProfile == "1" && (
                <div className={styles.profile}>
                  {profiles && media?.profile_id in profiles
                    ? profiles[media?.profile_id]?.name
                    : ""}
                </div>
              )}
            </div>
          </div>
        </div>
      </Link>
    </div>
  );
};
export default PosterComponent;
