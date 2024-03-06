import styles from "./MassEditor.module.scss";
import { useContext, useEffect, useState } from "react";
import { WebSocketContext } from "../../contexts/webSocketContext";
import { ReactComponent as BookmarkFilled } from "../svgs/bookmark_filled.svg";
import { ReactComponent as BookmarkUnfilled } from "../svgs/bookmark_unfilled.svg";
import { ReactComponent as ContinuingIcon } from "../svgs/play_arrow.svg";
import { ReactComponent as StoppedIcon } from "../svgs/stop.svg";
import InputCheckbox from "../inputs/inputCheckbox/InputCheckbox";
import InputSelect from "../inputs/inputSelect/InputSelect";
import MassEditorToolbar from "../toolbars/massEditorToolbar/MassEditorToolbar";
import sortAndFilter from "../../utils/sortAndFilter";
import { formatSize } from "../../utils/format";

const MassEditor = () => {
	const wsContext: any = useContext(WebSocketContext);
	const series: any = wsContext?.data?.series;
	const settings: any = wsContext?.data?.settings;
	const seriesArray = Array.from(Object.values(series || {}));
	const profiles: any = wsContext?.data?.profiles;
	const [selectedSeries, setSelectedSeries] = useState<any>([]);
	const [monitored, setMonitored] = useState<any>(false);
	const [profile, setProfile] = useState<any>();
	const [selected, setSelected] = useState<string | null>(null);
	const [selectAll, setSelectAll] = useState(false);

	const sort = settings?.massEditor_sort;
	const sortDirection = settings?.massEditor_sort_direction;
	const filter = settings?.massEditor_filter;
	const sortedSeries = sortAndFilter(
		series,
		profiles,
		sort,
		sortDirection,
		filter,
	);

	const applyChanges = () => {
		for (const series of selectedSeries) {
			series.monitored =
				parseInt(monitored) !== -1 ? parseInt(monitored) : undefined;
			series.profile_id =
				parseInt(profile) !== 0 ? parseInt(profile) : undefined;
			fetch(`http://${window.location.hostname}:7889/api/series/${series.id}`, {
				method: "PUT",
				headers: {
					"Content-Type": "application/json",
					Authorization: `Bearer ${localStorage.getItem("token")}`,
				},

				body: JSON.stringify(series),
			});
		}
	};

	const handleCheckboxChange = (series: any) => {
		setSelectedSeries((prevSelected: any[]) =>
			prevSelected.some((s) => s.id === series.id)
				? prevSelected.filter((s) => s.id !== series.id)
				: [...prevSelected, series],
		);
	};
	const handleSelectAllChange = () => {
		setSelectAll(!selectAll);
		setSelectedSeries(!selectAll ? seriesArray : []);
	};

	useEffect(() => {
		applyChanges();
	}, [monitored, profile]);

	return (
		<div className={styles.massEditor}>
			<MassEditorToolbar
				selected={selected}
				setSelected={setSelected}
				settings={settings}
			/>
			<div className={styles.content}>
				{sortedSeries && sortedSeries.length !== 0 ? (
					<>
						<table className={styles.table}>
							<thead>
								<tr className={styles.headRow}>
									<th>
										<InputCheckbox
											checked={selectAll}
											onChange={handleSelectAllChange}
										/>
									</th>
									<th></th>
									<th>Series</th>
									<th>Profile</th>
									<th>Path</th>
									<th>Space Saved</th>
									<th>Size on Disk</th>
								</tr>
							</thead>
							<tbody>
								{sortedSeries?.map((s: any, index: any) => (
									<tr className={styles.row} key={index}>
										<td className={styles.inputCell}>
											<InputCheckbox
												checked={selectedSeries.some(
													(series: any) => series.id === s.id,
												)}
												onChange={() => handleCheckboxChange(s)}
											/>
										</td>
										<td className={styles.iconCell}>
											{s?.monitored ? (
												<BookmarkFilled className={styles.monitored} />
											) : (
												<BookmarkUnfilled className={styles.monitored} />
											)}
											{s?.status !== "Ended" ? (
												<ContinuingIcon className={styles.continue} />
											) : (
												<StoppedIcon className={styles.stopped} />
											)}
										</td>
										<td>
											<a href={"/series/" + s?.id} className={styles.name}>
												{s?.id}
											</a>
										</td>
										<td>{profiles ? profiles[s.profile_id]?.name : ""}</td>
										<td>/series/{s.id}</td>
										<td>{formatSize(s.space_saved)}</td>
										<td>{formatSize(s.size)}</td>
									</tr>
								))}
							</tbody>
						</table>
					</>
				) : (
					<>No Media Found</>
				)}
			</div>
			<div className={styles.footer}>
				<div className={styles.input}>
					<div className={styles.inputContainer}>
						<label className={styles.label}>Monitored </label>
						<InputSelect
							selected={monitored}
							onChange={(e: any) => {
								setMonitored(e.target.value);
							}}
						>
							<option value={-1}>{"No Change"}</option>
							<option value={0}>{"Not Monitored"}</option>
							<option value={1}>{"Monitored"}</option>
						</InputSelect>
					</div>
					<div className={styles.inputContainer}>
						<label className={styles.label}>Profile </label>
						<InputSelect
							selected={profile}
							onChange={(e: any) => {
								setProfile(e.target.value);
							}}
						>
							<option value={0}>{"No Change"}</option>
							{Object.values(profiles || {}).map(
								(profile: any, index: number) => (
									<option value={profile.id} key={index}>
										{profile.name}
									</option>
								),
							)}
						</InputSelect>
					</div>
				</div>
			</div>
		</div>
	);
};
export default MassEditor;
