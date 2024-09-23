import styles from "./ToolBarItem.module.scss";
import CheckmarkIcon from "../svgs/check.svg?react";
import ArrowDropdownIcon from "../svgs/arrow_dropdown.svg?react";
import { useEffect, useRef } from "react";
const ToolBarOption = ({
  icon,
  index,
  text,
  settings,
  disabled,
  dropdownItems,
  selected,
  setSelected,
  sort = false,
  sortDirection,
  onClick = () => {
    undefined;
  },
}: any) => {
  const changeSelected = (index: number) => {
    onClick();
    if (selected === index) {
      setSelected(-1);
    } else {
      setSelected(index);
    }
  };
  const itemRef: any = useRef(null);
  useEffect(() => {
    function handleClickOutside(event: any) {
      if (itemRef.current && !itemRef.current.contains(event.target)) {
        setSelected(-1);
      }
    }

    if (selected === index) {
      document.addEventListener("mousedown", handleClickOutside);
    } else {
      document.removeEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [selected, index, setSelected]);
  return (
    <div className={styles.toolBarItem} ref={itemRef}>
      <div
        className={styles.toolBarOption}
        onClick={() => changeSelected(index)}
        style={
          disabled ? { opacity: "50%", cursor: "default", fill: "white" } : {}
        }
      >
        <div className={styles.svg}>{icon}</div>
        <div className={styles.text}>{text}</div>
      </div>
      {selected === index && (
        <div className={styles.dropdown}>
          {dropdownItems?.map((item: any, itemIndex: number) => (
            <div key={itemIndex}>
              <div
                className={styles.item}
                onClick={(e) => {
                  e.stopPropagation();
                  item.onClick();
                  setSelected(-1);
                }}
              >
                <div className={styles.text}>{item?.text}</div>
                {item?.id === settings[item?.settingId] || "" ? (
                  <div className={styles.svgContainer}>
                    {sort ? (
                      <ArrowDropdownIcon
                        style={
                          sortDirection == "descending"
                            ? { transform: "rotate(180deg)" }
                            : {}
                        }
                        className={styles.svg}
                      />
                    ) : (
                      <>
                      <CheckmarkIcon className={styles.svg} />
                      </>
                    )}
                  </div>
                ) : (
                  <></>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
export default ToolBarOption;
