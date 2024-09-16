import styles from "./Poster.module.scss";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
const PosterComponent = ({ media, posterWidth, posterHeight, profiles, settings }: any) => {
  console.log(media, 'hi')
  console.log(settings, 'set')
  const type = media?.episodeCount != undefined ? "series" : "movies";
  const progress = () => {
    if (media?.episodeCount == undefined) {
      return media?.missing ? "0%" : "100%";
    }
    return media?.episodeCount === 0
      ? "100%"
      : ((media?.episodeCount - media?.missingEpisodes) /
          media?.episodeCount || 0) *
          100 +
          "%";
  };
  const backgroundColor = () => {
    if (progress() === "100%") {
      if (media?.missingEpisodes == undefined) {
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
                      settings.mediaPosterDetailedProgressBar
                        ? "15px"
                        : "5px",
                  }}
                />
                {settings?.mediaPosterDetailedProgressBar && (
                  <div className={styles.detailText}>
                    {media?.episodeCount == undefined ? (
                      <>{media?.missing ? "0/1" : "1/1"}</>
                    ) : (
                      <>
                        {media?.episodeCount - media?.missingEpisodes}/
                        {media?.episodeCount}
                      </>
                    )}
                  </div>
                )}
              </div>
              {settings?.mediaPosterShowTitle && (
                <div className={styles.name}>
                  {media?.name ? media?.name : media?.id}
                </div>
              )}
              {settings?.mediaPosterShowMonitored && (
                <div className={styles.status}>
                  {media?.monitored ? "Monitored" : "Unmonitored"}
                </div>
              )}
              {settings?.mediaPosterShowProfile && (
                <div className={styles.profile}>
                  {profiles && media?.profileId in profiles
                    ? profiles[media?.profileId]?.name
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
