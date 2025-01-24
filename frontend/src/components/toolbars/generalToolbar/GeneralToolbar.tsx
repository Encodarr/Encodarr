import styles from "./GeneralToolbar.module.scss";
import ToolBar from "../../toolBar/ToolBar";
import ToolBarItem from "../../toolBarItem/ToolBarItem";
import SaveIcon from "../../svgs/save.svg?react";

const GeneralToolbar = ({
	selected,
	setSelected,
	currentSettings,
	setSettingsChanged,
	settingsChanged,
}: any) => {
	const handleSave = () => {
		if (!settingsChanged) {
			return;
		}
		for (const key in currentSettings) {
			if (key == "username" || key == "password") {
				fetch(`/api/user/` + key, {
					method: "PUT",
					headers: {
						"Content-Type": "application/json",
						Authorization: `Bearer ${localStorage.getItem("token")}`,
					},
					body: JSON.stringify({username: currentSettings['username'], password: currentSettings['password']}),
				});
			}
			else {
				fetch(`/api/settings/` + key, {
					method: "PUT",
					headers: {
						"Content-Type": "application/json",
						Authorization: `Bearer ${localStorage.getItem("token")}`,
					},
					body: JSON.stringify({id: key, value: currentSettings[key]}),
				});
			}
		}
		setSettingsChanged(false);
	};
	const leftToolBarItems: any = [
		<ToolBarItem
			text={!settingsChanged ? "No Changes" : "Save"}
			key="save"
			icon={<SaveIcon className={styles.svg} />}
			onClick={handleSave}
			disabled={!settingsChanged}
			selected={selected}
			setSelected={setSelected}
		/>,
	];
	return <ToolBar leftToolBarItems={leftToolBarItems} />;
};
export default GeneralToolbar;
