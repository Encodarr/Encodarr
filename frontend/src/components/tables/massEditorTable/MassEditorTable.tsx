import InputCheckbox from "../../inputs/inputCheckbox/InputCheckbox";
import Table from "../../table/Table";
import styles from "./MassEditorTable.module.scss";
import { formatSize } from "../../../utils/format";
import BookmarkFilled from "../../svgs/bookmark_filled.svg?react";
import BookmarkUnfilled from "../../svgs/bookmark_unfilled.svg?react";
import ContinuingIcon from "../../svgs/play_arrow.svg?react";
import StoppedIcon from "../../svgs/stop.svg?react";
import { Tooltip } from "react-tooltip";
import { Link } from "react-router-dom";
const MassEditorTable = ({
	sortedMedia,
	selectedMedia,
	selectAll,
	handleSelectAllChange,
	handleCheckboxChange,
	profiles,
}: any) => {
	return (
		<Table>
			<thead>
				<tr className={styles.headRow}>
					<th>
						<InputCheckbox
							checked={selectAll}
							onChange={handleSelectAllChange}
						/>
					</th>
					<th></th>
					<th>Title</th>
					<th>Profile</th>
					<th>Path</th>
					<th>Space Saved</th>
					<th>Size on Disk</th>
				</tr>
			</thead>
			<tbody>
				{sortedMedia?.map((media: any, index: any) => (
					<tr className={styles.row} key={index}>
						<td className={styles.inputCell}>
							<InputCheckbox
								checked={selectedMedia.some(
									(series: any) => series.id === media.id
								)}
								onChange={() => handleCheckboxChange(media)}
							/>
						</td>
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
							<Tooltip id="monitoredTooltip" place="top" content="Monitored" />
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
							<Link to={
								(media?.missing_episodes == undefined
									? "/movies/"
									: "/series/") + media?.id
							} className={styles.name}>
								{media?.id}
							</Link>
						</td>
						<td>{profiles ? profiles.find((profile: any) => profile.id === media.profileId)?.name : ""}</td>
						<td>
							/{media?.missing_episodes ? "movies" : "series"}/
							{media.id}
						</td>
						<td>{formatSize(media?.spaceSaved || 0)}</td>
						<td>{formatSize(media?.size || 0)}</td>
					</tr>
				))}
			</tbody>
		</Table>
	);
};
export default MassEditorTable;
