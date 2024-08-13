import styles from "./Posters.module.scss";
import PosterComponent from "../poster/Poster";

const Posters = ({ sortedMedia, settings }: any) => {
  const size = settings?.media_poster_posterSize;
  let posterWidth = "128px";
  let posterHeight = "260px";
  if (size === "small") {
    posterWidth = "115px";
    posterHeight = "234px";
  } else if (size === "medium") {
    posterWidth = "128px";
    posterHeight = "260px";
  } else if (size === "large") {
    posterWidth = "170px";
    posterHeight = "324px";
  }

  return (
    <div className={styles.postersView}>
      <div className={styles.content}>
        {sortedMedia.map((media: any) => (
          <PosterComponent
            media={media}
            key={media.id}
            posterWidth={posterWidth}
            posterHeight={posterHeight}
          />
        ))}
      </div>
    </div>
  );
};

export default Posters;
