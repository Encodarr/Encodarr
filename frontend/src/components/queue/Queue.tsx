import styles from "./Queue.module.scss";
import { useContext, useState } from "react";
import { WebSocketContext } from "../../contexts/webSocketContext";
import QueueToolbar from "../toolbars/queueToolbar/QueueToolbar";
import ResetWrench from "../svgs/reset_wrench.svg?react";
import QueueModal from "../modals/queueModal/QueueModal";
import { formatETA } from "../../utils/format";
import QueueTable from "../tables/queueTable/QueueTable";
import { Link } from "react-router-dom";

const Queue = () => {
	const wsContext = useContext(WebSocketContext);
	const profiles = wsContext?.data?.profiles;
	const series = wsContext?.data?.series;
	const queue = wsContext?.data?.queue;
	const settings: any = wsContext?.data?.settings
		? Object.keys(wsContext?.data?.settings).reduce((acc, key) => {
			acc[key] = wsContext?.data?.settings[key].value;
			return acc;
		}, {})
		: {};
	const [selected, setSelected] = useState<string | null>(null);
	const [isModalOpen, setIsModalOpen] = useState(false);
	const [content, setContent] = useState<any>({});
	const queueArray = Array.from(Object.values(queue?.queue || []));
	return (
		<div className={styles.queue}>
			<QueueToolbar
				settings={settings}
				setContent={setContent}
				setIsModalOpen={setIsModalOpen}
				selected={selected}
				setSelected={setSelected}
			/>
			<QueueModal
				isOpen={isModalOpen}
				setIsOpen={setIsModalOpen}
				content={content}
				setContent={setContent}
			/>
			<div className={styles.content}>
				<>
					<div className={styles.currentContainer}>
						<div className={styles.current}>
							<table className={styles.table}>
								<thead>
									<tr className={styles.headRow}>
										<th></th>
										<th>Title</th>
										<th>Type</th>
										<th>Episode</th>
										<th>Profile</th>
										<th>Stage</th>
										<th>ETA</th>
										<th>Progress</th>
									</tr>
								</thead>
								<tbody>
									<tr className={styles.row}>
										<td className={styles.iconCell}>
											<ResetWrench className={styles.svg} />
										</td>
										<td className={styles.name}>
											<Link to={
												queue.current?.type == "series"
													? "/series/" + queue.current?.SeriesId
													: "/movies/" + queue.current?.id
											} className={styles.name}>
												{queue.current?.type == "series"
													? queue.current?.seriesId
													: queue.current?.id}
											</Link>
										</td>
										<td>
											{queue?.current && (
												<>
													{queue?.current && queue?.current?.type == "series"
														? "Series"
														: "Movie"}
												</>
											)}
										</td>
										<td>
											{queue?.current && queue?.current?.seriesId ? (
												<>
													{queue?.current?.seasonNumber}x
													{queue?.current?.episodeNumber}
												</>
											) : (
												<></>
											)}
										</td>
										<td>
											{queue?.current?.seriesId ? (
												<>
													{profiles && series && (
														profiles.find((profile: any) =>
															profile.id === series.find((series: any) =>
																series.id == queue?.current?.seriesId
															)?.profileId
														)?.name ?? 'Profile not found'
													)}
												</>
											) : (
												<>
													{profiles ? profiles.find((profile: any) => profile.id === queue?.current?.profileId)?.name : ""}
												</>
											)}
										</td>
										<td>
											{settings?.queueStatus == "active"
												? queue?.stage || "--"
												: "Paused"}
										</td>
										<td>{queue?.current.id != "" && formatETA(queue?.eta)}</td>
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
														width: `${queue?.progress || 0}%`,
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
					</div>
					{queueArray?.length !== 0 ? (
						<>
							<QueueTable
								series={series}
								profiles={profiles}
								queueArray={queueArray}
								settings={settings}
							/>
						</>
					) : (
						<></>
					)}
				</>
			</div>
		</div>
	);
};
export default Queue;
