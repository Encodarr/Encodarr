import { useContext, useEffect, useRef, useState } from "react";
import styles from "./Movie.module.scss";
import Drive from "../svgs/hard_drive.svg?react";
import Profile from "../svgs/person.svg?react";
import Monitored from "../svgs/bookmark_filled.svg?react";
import Unmonitored from "../svgs/bookmark_unfilled.svg?react";
import Continuing from "../svgs/play_arrow.svg?react";
import Ended from "../svgs/stop.svg?react";
import Network from "../svgs/tower.svg?react";
import { WebSocketContext } from "../../contexts/webSocketContext";
import { formatSize } from "../../utils/format";
import FolderIcon from "../svgs/folder.svg?react";
import MovieToolbar from "../toolbars/movieToolbar/MovieToolbar";
import MovieModal from "../modals/movieModal/MovieModal";

const Movie = ({ movieName }: any) => {
	const wsContext = useContext(WebSocketContext);
	const profiles = wsContext?.data?.profiles;
	const movie: any = 
	wsContext?.data?.movies && profiles
	  ? wsContext?.data?.movies.find((s: any) => s.id === movieName)
	  : {};
	  const system: any = wsContext?.data?.system
	  ? Object.keys(wsContext?.data?.system).reduce((acc, key) => {
		acc[key] = wsContext?.data?.system[key].value;
		return acc;
	  }, {})
	  : {};
	const [content, setContent] = useState<any>({});
	const handleEditClick = () => {
		setIsModalOpen(true);
		setContent(movie);
	};
	const [selected, setSelected] = useState<string | null>(null);
	const [isModalOpen, setIsModalOpen] = useState(false);
	const status = movie?.status;
	const genre = movie?.genre;
	const firstAirDate = movie?.releaseDate?.split("-")[0].trim();
	const lastAirDate = movie?.lastAirDate?.split("-")[0].trim();
	const overview = movie?.overview;
	const runYears =
		status === "Ended" ? firstAirDate + "-" + lastAirDate : firstAirDate;
	const [backdropSrc, setBackdropSrc] = useState<string | null>("");
	const [posterSrc, setPosterSrc] = useState<string | null>("");
	const loaded = useRef(false);

	useEffect(() => {
		if (loaded.current == true) {
			return;
		}
		const fetchImage = async (
			path: string,
			setSrc: (src: string | null) => void
		) => {
			try {
				let cache = null;
				let cachedResponse = null;
				if ("caches" in window) {
					cache = await caches.open("image-cache");
					cachedResponse = await cache.match(
						`/api/movies/${movie?.id}/${path}`
					);
				}

				if (cachedResponse) {
					const blob = await cachedResponse.blob();
					setSrc(URL.createObjectURL(blob));
				} else {
					const response = await fetch(`/api/movies/${movie?.id}/${path}`, {
						headers: {
							Authorization: `Bearer ${localStorage.getItem("token")}`,
						},
					});
					if (response.status !== 200) {
						setSrc(null);
						return;
					}
					const clonedResponse = response.clone();
					const blob = await response.blob();
					setSrc(URL.createObjectURL(blob));
					if (cache) {
						cache.put(`/api/movies/${movie?.id}/${path}`, clonedResponse);
					}
				}
			} catch (e) {
				console.log(e);
			}
		};

		if (movie?.id && movie?.id !== "") {
			fetchImage("backdrop", setBackdropSrc);
			fetchImage("poster", setPosterSrc);
			loaded.current = true;
		}
	}, [movie?.id]);
	const profileName = profiles?.find((profile: any) => profile.id === movie?.profileId)?.name || "";

	return (
		<div className={styles.movie}>
			<MovieToolbar
				system={system}
				selected={selected}
				setSelected={setSelected}
				handleEditClick={handleEditClick}
				movieName={movieName}
			/>
			<MovieModal
				isOpen={isModalOpen}
				setIsOpen={setIsModalOpen}
				content={content}
				setContent={setContent}
				profiles={profiles}
			/>
			<div className={styles.movieContent}>
				<div className={styles.header}>
					<img
						className={styles.backdrop}
						src={backdropSrc || "/backdrop.jpg"}
						alt="backdrop"
					/>
					<div className={styles.filter} />
					<div className={styles.content}>
						<img
							className={styles.poster}
							src={posterSrc || "/poster.png"}
							alt="poster"
						/>
						<div className={styles.details}>
							<div className={styles.titleRow}>
								<div className={styles.headerIcon}>
									{movie?.monitored ? (
										<Monitored className={styles.monitoredSVG} />
									) : (
										<Unmonitored className={styles.monitoredSVG} />
									)}
								</div>

								{movie?.name ? movie?.name : movie?.id}
							</div>
							<div className={styles.movieDetails}>
								<span className={styles.runtime}>
									{movie?.runtime ? movie?.runtime : "-"} Minutes
								</span>
								{genre ? <span className={styles.genre}>{genre}</span> : <></>}
								{status ? (
									<span className={styles.runYears}>{runYears}</span>
								) : (
									<></>
								)}
							</div>
							<div className={styles.tags}>
								<div className={styles.tag}>
									<div className={styles.icon}>
										<FolderIcon className={styles.svg} />
									</div>
									{"/movies/" + movie?.id}
								</div>

								<div className={styles.tag}>
									<div className={styles.icon}>
										<Drive className={styles.svg} />
									</div>
									{formatSize(movie?.size)}
								</div>
								<div className={styles.tag}>
									<div className={styles.icon}>
										<Profile className={styles.svg} />
									</div>
									{profileName}
								</div>
								<div className={styles.tag}>
									<div className={styles.icon}>
										{movie?.monitored ? (
											<Monitored className={styles.svg} />
										) : (
											<Unmonitored className={styles.svg} />
										)}
									</div>
									{movie?.monitored ? "Monitored" : "Unmonitored"}
								</div>
								{status ? (
									<div className={styles.tag}>
										<div className={styles.icon}>
											{status === "Ended" ? (
												<Ended className={styles.svg} />
											) : (
												<Continuing className={styles.svg} />
											)}
										</div>
										{status}
									</div>
								) : (
									<></>
								)}
								{movie?.studio ? (
									<div className={styles.tag}>
										<div className={styles.icon}>
											<Network className={styles.svg} />
										</div>
										{movie?.studio}
									</div>
								) : (
									<></>
								)}
							</div>
							<div className={styles.overview}>{overview}</div>
						</div>
					</div>
				</div>
				<div className={styles.fileContainer}>
					<table>
						<thead>
							<tr>
								<th>Filename</th>
								<th>Video Codec</th>
								<th>Size</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td>{movie?.filename}</td>
								<td>{movie?.videoCodec}</td>
								<td>{formatSize(movie?.size)}</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
	);
};
export default Movie;
