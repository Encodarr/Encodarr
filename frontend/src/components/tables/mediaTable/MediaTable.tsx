import styles from "./MediaTable.module.scss";
import BookmarkFilled from "../../svgs/bookmark_filled.svg?react";
import BookmarkUnfilled from "../../svgs/bookmark_unfilled.svg?react";
import ContinuingIcon from "../../svgs/play_arrow.svg?react";
import StoppedIcon from "../../svgs/stop.svg?react";
import { Link } from "react-router-dom";
import { formatSize } from "../../../utils/format";
import Table from "../../table/Table";
import { Tooltip } from "react-tooltip";
const MediaTable = ({ settings, profiles, sortedMedia }: any) => {
	const progress = (series: any) => {
		return (
			((series?.episodeCount - series?.missingEpisodes) /
				series?.episodeCount || 0) *
				100 +
			"%"
		);
	};

	const backgroundColor = (series: any) => {
		if (progress(series) === "100%") {
			return series?.status === "Ended"
				? "rgb(39, 194, 76)"
				: "rgb(93, 156, 236)";
		} else {
			return series?.monitored ? "rgb(240, 80, 80)" : "rgb(255, 165, 0)";
		}
	};
	return (
		<div className={styles.tableContainer}>
			<Table>
				<thead>
					<tr>
						<th></th>
						<th>Title</th>
						{settings?.mediaTableShowType == "true" && <th>Type</th>}
						{settings?.mediaTableShowProfile == "true" && <th>Profile</th>}
						{settings?.mediaTableShowNetwork == "true" && (
							<th>Network / Studio</th>
						)}
						{settings?.mediaTableShowSeasons == "true" && <th>Seasons</th>}
						{settings?.mediaTableShowEpisodes == "true" && <th>Episodes</th>}
						{settings?.mediaTableShowEpisodeCount == "true" && (
							<th>Episode Count</th>
						)}
						{settings?.mediaTableShowYear == "true" && <th>Year</th>}
						{settings?.mediaTableShowPath == "true" && <th>Path</th>}
						{settings?.mediaTableShowSizeOnDisk == "true" && (
							<th>Size On Disk</th>
						)}
						{settings?.mediaTableShowSizeSaved == "true" && <th>Space Saved</th>}
						{settings?.mediaTableShowGenre == "true" && <th>Genre</th>}
					</tr>
				</thead>
				<tbody>
					{sortedMedia.map((media: any, key: number) => (
						<tr className={styles.row} key={key}>
							<td className={styles.iconCell}>
								{media?.monitored ? (
									<BookmarkFilled
										data-tooltip-id="monitoredTooltip"
										className={styles.svg}
									/>
								) : (
									<BookmarkUnfilled
										data-tooltip-id="unmonitoredTooltip"
										className={styles.svg}
									/>
								)}
								{media?.status !== "Ended" ? (
									<ContinuingIcon
										data-tooltip-id="continuingTooltip"
										className={styles.svg}
									/>
								) : (
									<StoppedIcon
										data-tooltip-id="stoppedTooltip"
										className={styles.svg}
									/>
								)}
								<Tooltip
									id="monitoredTooltip"
									place="top"
									content="Monitored"
								/>
								<Tooltip
									id="unmonitoredTooltip"
									place="top"
									content="Unmonitored"
								/>
								<Tooltip
									id="continuingTooltip"
									place="top"
									content="Continuing"
								/>
								<Tooltip id="stoppedTooltip" place="top" content="Stopped" />
							</td>
							<td>
								<Link to={"/series/" + media?.id} className={styles.name}>
									{media?.name ? media?.name : media?.id}
								</Link>
							</td>
							<td>{media?.episodeCount == undefined ? "Movie" : "Series"}</td>
							{settings?.mediaTableShowProfile == "true" && (
						<td>{profiles ? profiles.find((profile: any) => profile.id === media.profileId)?.name : ""}</td>
					)}
							{settings?.mediaTableShowNetwork == "true" && (
								<td>
									{media?.episodeCount == undefined
										? media?.studio
										: media?.networks}
								</td>
							)}
							{settings?.mediaTableShowSeasons == "true" && (
								<td>{media?.seasonsCount}</td>
							)}
							{settings?.mediaTableShowEpisodes == "true" && (
								<td>
									{media?.episodeCount != undefined && (
										<div className={styles.progressBar}>
											<div
												className={styles.progress}
												style={{
													backgroundColor: backgroundColor(media),
													width: progress(media),
												}}
											/>
											<div className={styles.detailText}>
												{media?.episodeCount - media?.missingEpisodes}/
												{media?.episodeCount}
											</div>
										</div>
									)}
								</td>
							)}
							{settings?.mediaTableShowEpisodeCount == "true" && (
								<td>{media?.episodeCount}</td>
							)}
							{settings?.mediaTableShowYear == "true" && (
								<td>{media?.releaseDate}</td>
							)}
							{settings?.mediaTableShowSpaceSaved == "true" && (
								<td>{formatSize(media.spaceSaved)}</td>
							)}
							{settings?.mediaTableShowSizeOnDisk == "true" && (
								<td>{formatSize(media.size)}</td>
							)}
							{settings?.mediaTableShowSizeSaved == "true" && (
								<td>{formatSize(media.spaceSaved)}</td>
							)}
							{settings?.mediaTableShowGenre == "true" && <td>{media.genre}</td>}
						</tr>
					))}
				</tbody>
			</Table>
		</div>
	);
};
export default MediaTable;
