import { useState } from "react";
import styles from "./Season.module.scss";
import Open from "../svgs/expand_circle_up.svg?react";
import Close from "../svgs/expand_circle_down.svg?react";
import MonitoredIcon from "../svgs/bookmark_filled.svg?react";
import UnmonitoredIcon from "../svgs/bookmark_unfilled.svg?react";
import { formatDate, formatSize } from "../../utils/format";
const Season = ({ season, monitored }: any) => {
	const [isOpen, setIsOpen] = useState(false);
	const onSeasonClick = () => {
		setIsOpen(!isOpen);
	};
	console.log(season)
	const backgroundColor = () => {
		if (season?.missingEpisodes != 0) {
			return "var(--dangerColor)";
		} else {
			return "var(--successColor)";
		}
	};
	return (
		<div className={styles.season}>
			<div className={styles.seasonHeader} onClick={onSeasonClick}>
				<div className={styles.left}>
					{monitored ? (
						<MonitoredIcon className={styles.svg} />
					) : (
						<UnmonitoredIcon className={styles.svg} />
					)}
					<div className={styles.seasonNumber}>{season?.name}</div>
					<div
						className={styles.profileRatio}
						style={{ backgroundColor: backgroundColor() }}
					>
						{season?.episodeCount - season?.missingEpisodes} /{" "}
						{season?.episodeCount}
					</div>
					<div className={styles.size}>{formatSize(season?.size)}</div>
				</div>
				<div className={styles.center}>
					<div className={styles.open}>
						{isOpen ? (
							<>
								<Open className={styles.openSVG} />
							</>
						) : (
							<Close className={styles.openSVG} />
						)}
					</div>
				</div>
			</div>
			<div
				className={styles.seasonInfo}
				style={isOpen ? {} : { display: "none" }}
			>
				<table>
					<thead>
						<tr>
							<th></th>
							<th>#</th>
							<th>Title</th>
							<th>Air Date</th>
							<th>Status</th>
							<th>Size</th>
						</tr>
					</thead>
					<tbody>
						{Object.values(season.episodes || {})
							.reverse()
							.map((episode: any, index: number) => (
								<tr key={index}>
									<td>{monitored ? <></> : <></>}</td>
									<td>{episode?.episodeNumber}</td>
									<td>
										{episode?.episodeName
											? episode.episodeName
											: episode.filename}
									</td>
									<td>{formatDate(episode?.airDate)}</td>
									<td>{episode?.videoCodec}</td>
									<td>{formatSize(episode?.size)}</td>
								</tr>
							))}
					</tbody>
				</table>
				<div className={styles.seasonInfoFooter} onClick={onSeasonClick}>
					<Open className={styles.svg} />
				</div>
			</div>
		</div>
	);
};
export default Season;
