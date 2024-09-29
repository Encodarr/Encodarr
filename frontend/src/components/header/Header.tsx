import styles from "./Header.module.scss";
import Person from "../svgs/person.svg?react";
import Logo from "../svgs/transfigurr.svg?react";
import { WebSocketContext } from "../../contexts/webSocketContext";
import { useContext, useEffect, useRef, useState } from "react";
import Timer from "../svgs/timer.svg?react";
import Build from "../svgs/build.svg?react";
import Pending from "../svgs/pending.svg?react";
import RestartIcon from "../svgs/restart.svg?react";
import ShutdownIcon from "../svgs/shutdown.svg?react";
import { Link } from "react-router-dom";
import { formatETA } from "../../utils/format";

const HeaderComponent = () => {
	const wsContext = useContext(WebSocketContext);
	const queue = wsContext?.data?.queue;
	const [openDropdown, setOpenDropdown] = useState(false);
	const dropdownRef: any = useRef(null);
	useEffect(() => {
		function handleClickOutside(event: any) {
			if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
				setOpenDropdown(false);
			}
		}

		if (openDropdown) {
			document.addEventListener("mousedown", handleClickOutside);
		} else {
			document.removeEventListener("mousedown", handleClickOutside);
		}

		return () => {
			document.removeEventListener("mousedown", handleClickOutside);
		};
	}, [openDropdown]);

	const shutdown = async () => {
		setOpenDropdown(false);
		await fetch(`/api/actions/shutdown`, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${localStorage.getItem("token")}`,
			},
		});
	};

	const restart = async () => {
		setOpenDropdown(false);

		await fetch(`/api/actions/restart`, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${localStorage.getItem("token")}`,
			},
		});
	};

	const settings: any = wsContext?.data?.settings
	? Object.keys(wsContext?.data?.settings).reduce((acc, key) => {
		acc[key] = wsContext?.data?.settings[key].value;
		return acc;
	  }, {})
	: {};
	return (
		<div className={styles.header}>
			<div className={styles.left}>
				<Link to="/" className={styles.logo}>
					<Logo className={styles.logoSVG} />
				</Link>
				<div className={styles.status}>
					<div className={styles.line}>
						<div className={styles.icon}>
							<Timer className={styles.svg} />
						</div>
						<div className={styles.text}>
							{queue && queue.stage !== "Idle"
								? Math.floor(queue?.progress)
								: "--"}
							%
						</div>
					</div>
					<div className={styles.line}>
						<div className={styles.icon}>
							<Pending className={styles.svg} />
						</div>
						<div className={styles.text}>
							{queue && queue.stage !== "Idle" ? formatETA(queue?.eta) : "--"}
						</div>
					</div>
					<div className={styles.line}>
						<div className={styles.icon}>
							<Build className={styles.svg} />
						</div>

						<div className={styles.text}>
							{settings?.queueStatus == "active"
								? queue?.stage || "--"
								: "Paused"}
						</div>
					</div>
				</div>
			</div>
			<div className={styles.right}>
				<div className={styles.profile} ref={dropdownRef}>
					<div
						className={styles.svg}
						onClick={() => setOpenDropdown(!openDropdown)}
					>
						<Person className={styles.svg} />
					</div>

					{openDropdown ? (
						<div className={styles.dropdown}>
							<div className={styles.item} onClick={restart}>
								<div className={styles.profilesvg}>
									<RestartIcon className={styles.actionSVG} />
								</div>
								<div className={styles.text}>Restart</div>
							</div>
							<div className={styles.item} onClick={shutdown}>
								<div className={styles.profilesvg}>
									<ShutdownIcon className={styles.actionSVG} />
								</div>
								<div className={styles.text}>Shutdown</div>
							</div>
						</div>
					) : (
						<></>
					)}
				</div>
			</div>
		</div>
	);
};
export default HeaderComponent;
