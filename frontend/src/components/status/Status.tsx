import styles from "./Status.module.scss";
import { useContext, useEffect, useState } from "react";
import { WebSocketContext } from "../../contexts/webSocketContext";
import { formatSize } from "../../utils/format";
import packageJson from "../../../package.json";
const version = packageJson.version;
const Status = () => {
	const wsContext = useContext(WebSocketContext);
	const system: any = wsContext?.data?.system;

	const [delta, setDelta] = useState("");

	useEffect(() => {
		const calculateDelta = () => {
			const timestamp = new Date(system?.start_time);
			const now = new Date();
			const delta = Math.abs(now.getTime() - timestamp.getTime());

			const seconds = Math.floor((delta / 1000) % 60);
			const minutes = Math.floor((delta / 1000 / 60) % 60);
			const hours = Math.floor((delta / (1000 * 60 * 60)) % 24);
			const days = Math.floor(delta / (1000 * 60 * 60 * 24));

			setDelta(
				days +
					"D " +
					hours.toString().padStart(2, "0") +
					":" +
					minutes.toString().padStart(2, "0") +
					":" +
					seconds.toString().padStart(2, "0")
			);
		};

		calculateDelta();
		const intervalId = setInterval(calculateDelta, 1000);

		return () => clearInterval(intervalId);
	}, [system?.start_time]);

	return (
		<div className={styles.status}>
			<div className={styles.diskSpace}>
				<div className={styles.header}>Disk Space</div>
				<table className={styles.table}>
					<thead>
						<tr className={styles.headRow}>
							<th>Location</th>
							<th>Free Space</th>
							<th>Total Space</th>
							<th>Bar</th>
						</tr>
					</thead>
					<tbody>
						<tr className={styles.row}>
							<td>/config</td>
							<td>{formatSize(system?.config_free_space)}</td>
							<td>{formatSize(system?.config_total_space)}</td>
							<td>
								<div
									style={{
										height: "20px",
										width: "100%",
										backgroundColor: "var(--progressBarBackgroundColor)",
										boxShadow: "inset 0 1px 2px rgba(0, 0, 0, 0.1)",
										borderRadius: "4px",
									}}
								>
									<div
										style={{
											height: "100%",
											width: `${((system?.config_total_space -
												system?.config_free_space) /
												system?.config_total_space) *
												100}%`,
											backgroundColor: "var(--transfigurrPurple)",
											borderRadius: "4px",
										}}
									/>
								</div>
							</td>
						</tr>
						<tr className={styles.row}>
							<td>/movies</td>
							<td>{formatSize(system?.movies_free_space)}</td>
							<td>{formatSize(system?.movies_total_space)}</td>
							<td>
								<div
									style={{
										height: "20px",
										width: "100%",
										backgroundColor: "var(--progressBarBackgroundColor)",
										boxShadow: "inset 0 1px 2px rgba(0, 0, 0, 0.1)",
										borderRadius: "4px",
									}}
								>
									<div
										style={{
											height: "100%",
											width: `${((system?.movies_total_space -
												system?.movies_free_space) /
												system?.movies_total_space) *
												100}%`,
											backgroundColor: "var(--transfigurrPurple)",
											borderRadius: "4px",
										}}
									/>
								</div>
							</td>
						</tr>
						<tr className={styles.row}>
							<td>/series</td>
							<td>{formatSize(system?.series_free_space)}</td>
							<td>{formatSize(system?.series_total_space)}</td>
							<td>
								<div
									style={{
										height: "20px",
										width: "100%",
										backgroundColor: "var(--progressBarBackgroundColor)",
										boxShadow: "inset 0 1px 2px rgba(0, 0, 0, 0.1)",
										borderRadius: "4px",
									}}
								>
									<div
										style={{
											height: "100%",
											width: `${((system?.series_total_space -
												system?.series_free_space) /
												system?.series_total_space) *
												100}%`,
											backgroundColor: "var(--transfigurrPurple)",
											borderRadius: "4px",
										}}
									/>
								</div>
							</td>
						</tr>
						<tr className={styles.row}>
							<td>/transcode</td>
							<td>{formatSize(system?.transcode_free_space)}</td>
							<td>{formatSize(system?.transcode_total_space)}</td>
							<td>
								<div
									style={{
										height: "20px",
										width: "100%",
										backgroundColor: "var(--progressBarBackgroundColor)",
										boxShadow: "inset 0 1px 2px rgba(0, 0, 0, 0.1)",
										borderRadius: "4px",
									}}
								>
									<div
										style={{
											height: "100%",
											width: `${((system?.transcode_total_space -
												system?.transcode_free_space) /
												system?.transcode_total_space) *
												100}%`,
											backgroundColor: "var(--transfigurrPurple)",
											borderRadius: "4px",
										}}
									/>
								</div>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
			<div className={styles.about}>
				<div className={styles.header}>About</div>
				<div className={styles.row}>
					<div className={styles.left}>Version</div>
					<div className={styles.right}>{version}</div>
				</div>
				<div className={styles.row}>
					<div className={styles.left}>Appdata Directory</div>
					<div className={styles.right}>/config</div>
				</div>
				<div className={styles.row}>
					<div className={styles.left}>Uptime</div>
					<div className={styles.right}>{delta}</div>
				</div>
			</div>
			<div className={styles.moreInfo}>
				<div className={styles.header}>More Info</div>
				<div className={styles.row}>
					<div className={styles.left}>Homepage</div>
					<div className={styles.right}>
						<a href={"https://transfigurr.media"}>transfigurr.media</a>
					</div>
				</div>
				<div className={styles.row}>
					<div className={styles.left}>Source</div>
					<div className={styles.right}>
						<a href={"https://github.com/Transfigurr/Transfigurr"}>
							github.com/Transfigurr/Transfigurr
						</a>
					</div>
				</div>
				<div className={styles.row}>
					<div className={styles.left}>Feature Requests</div>
					<div className={styles.right}>
						<a href={"https://github.com/Transfigurr/Transfigurr/issues"}>
							github.com/Transfigurr/Transfigurr/issues
						</a>
					</div>
				</div>
			</div>
		</div>
	);
};
export default Status;
