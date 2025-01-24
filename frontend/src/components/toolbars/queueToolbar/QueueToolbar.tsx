import styles from "./QueueToolbar.module.scss";
import ToolBar from "../../toolBar/ToolBar";
import ToolBarItem from "../../toolBarItem/ToolBarItem";
import Pause from "../../svgs/pause_circle.svg?react";
import Start from "../../svgs/play_circle.svg?react";
import TableIcon from "../../svgs/table.svg?react";

const QueueToolbar = ({
  settings,
  setContent,
  setIsModalOpen,
  selected,
  setSelected,
}: any) => {
  const handleOptionsClick = () => {
    setContent(settings);
    setIsModalOpen(true);
  };
  const setSetting = async (key: string, value: any) => {
    await fetch(`/api/settings/${key}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify({ id: key, value: value }),
    });
  };
  const middleToolBarItems: any = [
    <ToolBarItem
      text="Status"
      index={4}
      key={4}
      settings={settings}
      icon={
        settings?.queueStatus === "inactive" ? (
          <Pause className={styles.svg} />
        ) : (
          <Start className={styles.svg} />
        )
      }
      selected={selected}
      setSelected={setSelected}
      dropdownItems={[
        {
          text: "Active",
          id: "active",
          key: "active",
          settingId: "queueStatus",
          onClick: () => setSetting("queueStatus", "active"),
        },
        {
          text: "Inactive",
          id: "inactive",
          key: "inactive",
          settingId: "queueStatus",
          onClick: () => setSetting("queueStatus", "inactive"),
        },
      ]}
    />,
  ];
  const rightToolbarItems: any = [
    <ToolBarItem
      text="Options"
      index={2}
      key={2}
      settings={settings}
      icon={<TableIcon className={styles.svg} />}
      onClick={handleOptionsClick}
      selected={selected}
      setSelected={setSelected}
    />,
  ];
  return (
    <ToolBar
      leftToolBarItems={[]}
      middleToolBarItems={middleToolBarItems}
      rightToolBarItems={rightToolbarItems}
    />
  );
};

export default QueueToolbar;
