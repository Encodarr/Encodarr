import styles from "./MovieToolbar.module.scss";
import ToolBarItem from "../../toolBarItem/ToolBarItem";
import ToolBar from "../../toolBar/ToolBar";
import SyncIcon from "../../svgs/cached.svg?react";
import EditIcon from "../../svgs/edit.svg?react";
import LoadingIcon from "../../svgs/loading.svg?react";
import RssFeedIcon from "../../svgs/rss_feed.svg?react";
const MovieToolbar = ({
	system,
	selected,
	setSelected,
	handleEditClick,
	movieName,
}: any) => {
	const handleScanClick = async () => {
		await fetch(`/api/actions/scan/movies/${movieName}`, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${localStorage.getItem("token")}`,
			},
		});
	};

	const handleMetadataClick = async () => {
		await fetch(`/api/actions/refresh/metadata/movies/${movieName}`, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${localStorage.getItem("token")}`,
			},
		});
	};
	const leftToolBarItems: any = [
		<ToolBarItem
			text="Scan"
			key="scan"
			icon={
				<SyncIcon
					className={
						system?.scanRunning == "1" && system?.scanTarget == movieName
							? styles.spinning
							: styles.svg
					}
				/>
			}
			onClick={handleScanClick}
			selected={selected}
			setSelected={setSelected}
		/>,
		<ToolBarItem
			text="Refresh Metadata"
			key="metadata"
			icon={
				system?.metadataRunning == "1" &&
				system?.metadataTarget == movieName ? (
					<LoadingIcon className={styles.loading} />
				) : (
					<RssFeedIcon className={styles.svg} />
				)
			}
			onClick={handleMetadataClick}
			selected={selected}
			setSelected={setSelected}
		/>,
		<ToolBarItem
			text="Edit"
			key="edit"
			icon={<EditIcon className={styles.svg} />}
			onClick={handleEditClick}
			selected={selected}
			setSelected={setSelected}
		/>,
	];
	return <ToolBar leftToolBarItems={leftToolBarItems} />;
};
export default MovieToolbar;
