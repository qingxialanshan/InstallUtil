package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ADT struct {
	Component
}

func (adt *ADT) Install(args ...string) {
	adt.ComponentId = args[0]
	adt.Version = args[1]
	adt.ExecFile = args[2]
	adt.InstallLocation = args[3]

	var elocation, nvpref string
	elocation = filepath.Join(args[3], "eclipse")
	if myutils.Global_OS == "windows" {
		elocation = filepath.Join(elocation, "eclipsec.exe")
		adt.ExecFile = "\"" + adt.ExecFile + "\""
	} else {
		elocation = filepath.Join(elocation, "eclipse")
		adt.ExecFile = strings.Replace(strings.Replace(strings.Replace(strings.Replace(adt.ExecFile, "\\", "\\\\", -1), "\"", "\\\"", -1), "'", "\\'", -1), " ", "\\ ", -1)
	}

	if myutils.Global_OS == "macosx" {
		nvpref = filepath.Join(args[3], "eclipse", "Eclipse.app", "Contents", "MacOS", "nvpref.ini")
	} else {
		nvpref = filepath.Join(args[3], "eclipse", "nvpref.ini")
	}

	//installing adt
	eargs := elocation + " -nosplash -application org.eclipse.equinox.p2.director -repository jar:file:" + adt.ExecFile + "!/, -installIU com.android.ide.eclipse.ddms.feature.feature.group,com.android.ide.eclipse.traceview.feature.feature.group,com.android.ide.eclipse.hierarchyviewer.feature.feature.group,com.android.ide.eclipse.adt.feature.feature.group,com.android.ide.eclipse.ndk.feature.feature.group,com.android.ide.eclipse.gldebugger.feature.feature.group"

	//fmt.Println(eargs)
	if myutils.Global_OS == "windows" {
		_, e := exec.Command("cmd", "/c", eargs).Output()
		myutils.CheckError(e)
	} else {
		_, e := exec.Command("bash", "-c", eargs).Output()
		myutils.CheckError(e)
	}

	//write to nvpref.ini
	sdk_path := filepath.Join(args[3], "android-sdk-"+myutils.Global_OS)
	config_str := `com.android.ide.eclipse.adt/com.android.ide.eclipse.adt.sdk=` + strings.Replace(sdk_path, "\\", "\\\\", -1) + `
org.eclipse.jdt.core/org.eclipse.jdt.core.classpathVariable.ANDROID_JAR=` + strings.Replace(sdk_path, "\\", "/", -1) + `/platforms/android-8/android.jar
com.android.ide.eclipse.adt/com.android.ide.eclipse.adt.skipPostCompileOnFileSave=false
org.eclipse.cdt.debug.mi.core/org.eclipse.cdt.debug.mi.core.SharedLibraries.auto_refresh=false
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.AmbiguousProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.RedefinitionProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.UnusedStaticFunctionProblem=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.FunctionResolutionProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.SuspiciousSemicolonProblem=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.ReturnStyleProblem=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.AssignmentInConditionProblem.params={}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.VariableResolutionProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.AssignmentToItselfProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.SuggestedParenthesisProblem=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.checkers.errnoreturn.params={implicit\=>;false}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.checkers.errreturnvalue.params={}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.AmbiguousProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.InvalidTemplateArgumentsProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.AssignmentToItselfProblem.params={}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.CatchByReference=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.CaseBreakProblem.params={no_break_comment\=>;"no break",last_case_param\=>;true,empty_case_param\=>;false}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.TypeResolutionProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.checkers.noreturn=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.StatementHasNoEffectProblem=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.FieldResolutionProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.MethodResolutionProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.checkers.errreturnvalue=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.RedeclarationProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.ScanfFormatStringSecurityProblem.params={}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.CatchByReference.params={unknown\=>;false,exceptions\=>;()}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.FieldResolutionProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.ReturnStyleProblem.params={}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.NonVirtualDestructorProblem.params={}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.checkers.errnoreturn=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.UnusedFunctionDeclarationProblem.params={macro\=>;true}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.RedeclarationProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.CircularReferenceProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.StatementHasNoEffectProblem.params={macro\=>;true,exceptions\=>;()}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.InvalidArguments=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.AbstractClassCreation.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.OverloadProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.AbstractClassCreation=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.UnusedFunctionDeclarationProblem=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.UnusedVariableDeclarationProblem=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.LabelStatementNotFoundProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.UnusedStaticFunctionProblem.params={macro\=>;true}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.AssignmentInConditionProblem=-Warning
org.eclipse.cdt.codan.core/eclipse.preferences.version=1
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.MemberDeclarationNotFoundProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.CircularReferenceProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.RedefinitionProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.NonVirtualDestructorProblem=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.SuspiciousSemicolonProblem.params={else\=>;false,afterelse\=>;false}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.NamingConventionFunctionChecker=-Info
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.InvalidTemplateArgumentsProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.InvalidArguments.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.SuggestedParenthesisProblem.params={paramNot\=>;false}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.TypeResolutionProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.FunctionResolutionProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.UnusedVariableDeclarationProblem.params={macro\=>;true,exceptions\=>;("@(\#)","$Id")}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.MethodResolutionProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.OverloadProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.LabelStatementNotFoundProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.VariableResolutionProblem=-Error
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.NamingConventionFunctionChecker.params={pattern\=>;"^[a-z]",macro\=>;true,exceptions\=>;()}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.MemberDeclarationNotFoundProblem.params={launchModes\=>;{RUN_ON_FULL_BUILD\=>;false,RUN_ON_INC_BUILD\=>;false,RUN_AS_YOU_TYPE\=>;true,RUN_ON_DEMAND\=>;true}}
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.CaseBreakProblem=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.internal.checkers.ScanfFormatStringSecurityProblem=-Warning
org.eclipse.cdt.codan.core/org.eclipse.cdt.codan.checkers.noreturn.params={implicit\=>;false}`
	_, fe := os.Stat(nvpref)
	if fe != nil {
		os.Create(nvpref)
	}
	myutils.Write_To_File(config_str, nvpref)
	fmt.Println(filepath.Join(adt.InstallLocation, "eclipse"))
}

func (adt *ADT) Uninstall(args ...string) {

	return
}
