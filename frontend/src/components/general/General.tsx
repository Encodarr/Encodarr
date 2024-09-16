import { useContext, useEffect, useRef, useState } from "react";
import styles from "./General.module.scss";
import { WebSocketContext } from "../../contexts/webSocketContext";
import InputContainer from "../inputs/inputContainer/InputContainer";
import GeneralToolbar from "../toolbars/generalToolbar/GeneralToolbar";

const General = () => {
	const wsContext = useContext(WebSocketContext);
	const initialSettings = wsContext?.data?.settings
  ? Object.keys(wsContext?.data?.settings).reduce((acc, key) => {
      acc[key] = wsContext?.data?.settings[key].value;
      return acc;
    }, {})
  : {};
	const profiles = wsContext?.data?.profiles;
	const [currentSettings, setCurrentSettings] = useState<any>({});
	const hasSetCurrentSettings = useRef(false);
	const [settingsChanged, setSettingsChanged] = useState(false);
	const [selected, setSelected] = useState<string | null>(null);
	const handleChange = (key: string, value: any) => {
		setCurrentSettings((prevSettings: any) => {
			const newSettings = {
				...prevSettings,
				[key]: value,
			};

			setSettingsChanged(
				JSON.stringify(newSettings) !== JSON.stringify(initialSettings)
			);
			return newSettings;
		});
	};
	useEffect(() => {
		if (initialSettings !== undefined && !hasSetCurrentSettings.current) {
			setCurrentSettings({
				...initialSettings,
				password: "passwordplaceholder",
			});
			hasSetCurrentSettings.current = true;
		}
	}, [initialSettings]);
	return (
		<div className={styles.general}>
			<GeneralToolbar
				selected={selected}
				setSelected={setSelected}
				settingsChanged={settingsChanged}
				setSettingsChanged={setSettingsChanged}
				currentSettings={currentSettings}
			/>
			<div className={styles.content}>
				<div className={styles.section}>
					<label className={styles.header}>Default Behavior</label>
					<InputContainer
						type="select"
						label="Default Profile"
						selected={currentSettings?.defaultProfile}
						onChange={(e: any) =>
							handleChange("defaultProfile", e.target.value)
						}
					>
						{Object.entries(profiles || {})?.map(([key, profile]: any) => (
							<option key={key} value={key}>
								{profile?.name}
							</option>
						))}
					</InputContainer>

					<InputContainer
						type="select"
						label="Queue State"
						selected={currentSettings?.queueStatus}
						onChange={(e: any) => handleChange("queueStatus", e.target.value)}
					>
						<option value="active">Active</option>
						<option value="inactive">Inactive</option>
					</InputContainer>

					<InputContainer
						type="select"
						label="Queue Startup State"
						selected={currentSettings?.queueStartupState}
						onChange={(e: any) =>
							handleChange("queueStartupState", e.target.value)
						}
					>
						<option value="previous">Previous</option>
						<option value="inactive">Inactive</option>
						<option value="active">Active</option>
					</InputContainer>
				</div>
				<div className={styles.section}>
					<label className={styles.header}>User Interface</label>
					<InputContainer
						type="select"
						label="Theme"
						selected={currentSettings?.theme}
						onChange={(e: any) => handleChange("theme", e.target.value)}
					>
						<option value="auto">Auto</option>
						<option value="light">Light</option>
						<option value="dark">Dark</option>
					</InputContainer>
				</div>
				<div className={styles.section}>
					<label className={styles.header}>Security</label>
					<InputContainer
						type="text"
						label="Username"
						warningText="Requires restart to take effect"
						selected={currentSettings?.username}
						onChange={(e: any) => handleChange("username", e.target.value)}
					/>
					<InputContainer
						type="password"
						label="Password"
						warningText="Requires restart to take effect"
						selected={currentSettings?.password}
						onChange={(e: any) => handleChange("password", e.target.value)}
					/>
				</div>
				<div className={styles.section}>
					<label className={styles.header}>Host</label>
					<InputContainer
						type="text"
						label="Internal Port"
						warningText="Requires restart to take effect"
						selected={currentSettings?.port}
						onChange={(e: any) => handleChange("port", e.target.value)}
					/>
				</div>
				<div className={styles.section}>
					<label className={styles.header}>Logging</label>
					<InputContainer
						type="select"
						label="Log Level"
						warningText="Requires restart to take effect"
						selected={currentSettings?.logLevel}
						onChange={(e: any) => handleChange("logLevel", e.target.value)}
					>
						<option value="info">Info</option>
						<option value="debug">Debug</option>
					</InputContainer>
				</div>
			</div>
		</div>
	);
};
export default General;
